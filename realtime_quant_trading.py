"""
实时股票量化交易示例（教育用途）
- 实时行情：Alpaca StockDataStream
- 策略：短周期/长周期均线交叉（基于最新价格流）
- 下单：Alpaca TradingClient（建议 paper 账户）

依赖：
    pip install alpaca-py pandas

环境变量：
    export APCA_API_KEY_ID='your_key'
    export APCA_API_SECRET_KEY='your_secret'
    export APCA_PAPER='true'   # true=模拟盘, false=实盘

运行：
    python realtime_quant_trading.py
"""

from __future__ import annotations

import asyncio
import os
from collections import deque
from dataclasses import dataclass

import pandas as pd
from alpaca.data.live.stock import StockDataStream
from alpaca.trading.client import TradingClient
from alpaca.trading.enums import OrderSide, TimeInForce
from alpaca.trading.requests import MarketOrderRequest


@dataclass
class StrategyConfig:
    symbol: str = "AAPL"
    short_window: int = 20
    long_window: int = 60
    trade_qty: int = 1


class RealtimeMACrossTrader:
    def __init__(self, config: StrategyConfig) -> None:
        self.config = config
        self.price_buffer = deque(maxlen=config.long_window)

        self.api_key = os.getenv("APCA_API_KEY_ID", "")
        self.api_secret = os.getenv("APCA_API_SECRET_KEY", "")
        self.paper = os.getenv("APCA_PAPER", "true").lower() == "true"

        if not self.api_key or not self.api_secret:
            raise RuntimeError("请先设置 APCA_API_KEY_ID / APCA_API_SECRET_KEY 环境变量")

        self.trading_client = TradingClient(
            api_key=self.api_key,
            secret_key=self.api_secret,
            paper=self.paper,
        )
        self.stream = StockDataStream(
            api_key=self.api_key,
            secret_key=self.api_secret,
        )

    def _get_position_qty(self) -> int:
        try:
            pos = self.trading_client.get_open_position(self.config.symbol)
            return int(float(pos.qty))
        except Exception:
            return 0

    def _place_market_order(self, side: OrderSide, qty: int) -> None:
        if qty <= 0:
            return

        order = MarketOrderRequest(
            symbol=self.config.symbol,
            qty=qty,
            side=side,
            time_in_force=TimeInForce.DAY,
        )
        submitted = self.trading_client.submit_order(order_data=order)
        print(f"[ORDER] {side.value.upper()} {qty} {self.config.symbol} -> {submitted.id}")

    def _signal(self) -> str | None:
        if len(self.price_buffer) < self.config.long_window:
            return None

        series = pd.Series(self.price_buffer)
        short_ma = series.tail(self.config.short_window).mean()
        long_ma = series.tail(self.config.long_window).mean()

        if short_ma > long_ma:
            return "BUY"
        if short_ma < long_ma:
            return "SELL"
        return None

    async def on_trade(self, trade) -> None:
        price = float(trade.price)
        self.price_buffer.append(price)
        print(f"[TICK] {trade.symbol} price={price:.4f} time={trade.timestamp}")

        signal = self._signal()
        if signal is None:
            return

        pos_qty = self._get_position_qty()

        if signal == "BUY" and pos_qty <= 0:
            self._place_market_order(OrderSide.BUY, self.config.trade_qty)

        elif signal == "SELL" and pos_qty > 0:
            self._place_market_order(OrderSide.SELL, pos_qty)

    async def run(self) -> None:
        account = self.trading_client.get_account()
        print(
            f"账户状态: {account.status}, 模式={'PAPER' if self.paper else 'LIVE'}, "
            f"交易标的: {self.config.symbol}"
        )

        self.stream.subscribe_trades(self.on_trade, self.config.symbol)
        print("开始接收实时行情 ... Ctrl+C 退出")
        self.stream.run()


async def main() -> None:
    config = StrategyConfig(
        symbol="AAPL",        # 可改为其他美股代码
        short_window=20,
        long_window=60,
        trade_qty=1,
    )
    trader = RealtimeMACrossTrader(config)
    await trader.run()


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("\n已退出。")
