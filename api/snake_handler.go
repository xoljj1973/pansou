package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SnakeGameHandler renders a self-contained Snake mini game.
func SnakeGameHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(snakeGameHTML))
}

const snakeGameHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>贪吃蛇小游戏</title>
  <style>
    :root {
      color-scheme: dark;
      --bg: #07111f;
      --panel: rgba(14, 24, 42, 0.82);
      --panel-strong: rgba(21, 34, 57, 0.95);
      --line: rgba(148, 163, 184, 0.25);
      --accent: #38bdf8;
      --snake: #22c55e;
      --snake-head: #86efac;
      --food: #fb7185;
      --text: #e5f0ff;
      --muted: #9fb2c8;
    }

    * { box-sizing: border-box; }

    body {
      margin: 0;
      min-height: 100vh;
      display: grid;
      place-items: center;
      padding: 24px;
      font-family: Inter, "PingFang SC", "Microsoft YaHei", system-ui, sans-serif;
      color: var(--text);
      background:
        radial-gradient(circle at 20% 10%, rgba(56, 189, 248, 0.26), transparent 28%),
        radial-gradient(circle at 80% 0%, rgba(34, 197, 94, 0.18), transparent 32%),
        linear-gradient(135deg, #020617, var(--bg));
    }

    .game-shell {
      width: min(94vw, 860px);
      display: grid;
      grid-template-columns: minmax(320px, 560px) 1fr;
      gap: 20px;
      align-items: stretch;
    }

    .board-card, .side-card {
      border: 1px solid var(--line);
      border-radius: 28px;
      background: var(--panel);
      box-shadow: 0 24px 90px rgba(0, 0, 0, 0.38);
      backdrop-filter: blur(14px);
    }

    .board-card {
      position: relative;
      padding: 18px;
    }

    canvas {
      width: 100%;
      aspect-ratio: 1 / 1;
      display: block;
      border-radius: 20px;
      background: #020617;
      box-shadow: inset 0 0 0 1px rgba(148, 163, 184, 0.2);
    }

    .side-card {
      padding: 22px;
      display: flex;
      flex-direction: column;
      gap: 18px;
    }

    h1 {
      margin: 0;
      font-size: clamp(28px, 5vw, 44px);
      letter-spacing: -0.06em;
    }

    .subtitle {
      margin: -8px 0 0;
      color: var(--muted);
      line-height: 1.6;
    }

    .stats {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 10px;
    }

    .stat {
      padding: 14px;
      border: 1px solid var(--line);
      border-radius: 18px;
      background: var(--panel-strong);
    }

    .label {
      display: block;
      color: var(--muted);
      font-size: 13px;
    }

    .value {
      display: block;
      margin-top: 6px;
      font-size: 30px;
      font-weight: 800;
      line-height: 1;
    }

    .controls {
      display: grid;
      gap: 10px;
    }

    button {
      border: 0;
      border-radius: 16px;
      padding: 13px 16px;
      cursor: pointer;
      color: #04111f;
      background: linear-gradient(135deg, var(--snake-head), var(--accent));
      font: inherit;
      font-weight: 800;
      box-shadow: 0 14px 30px rgba(56, 189, 248, 0.22);
    }

    button.secondary {
      color: var(--text);
      background: rgba(148, 163, 184, 0.14);
      box-shadow: none;
      border: 1px solid var(--line);
    }

    .tips {
      margin: 0;
      padding-left: 20px;
      color: var(--muted);
      line-height: 1.8;
    }

    .mobile-pad {
      display: none;
      grid-template-columns: repeat(3, 58px);
      justify-content: center;
      gap: 8px;
      margin-top: auto;
    }

    .mobile-pad button {
      padding: 14px 0;
      color: var(--text);
      background: rgba(148, 163, 184, 0.14);
      box-shadow: none;
      border: 1px solid var(--line);
    }

    .mobile-pad .up { grid-column: 2; }
    .mobile-pad .left { grid-column: 1; }
    .mobile-pad .down { grid-column: 2; }
    .mobile-pad .right { grid-column: 3; }

    .overlay {
      position: absolute;
      inset: 18px;
      display: grid;
      place-items: center;
      border-radius: 20px;
      background: rgba(2, 6, 23, 0.72);
      text-align: center;
      padding: 28px;
    }

    .overlay[hidden] { display: none; }

    .overlay h2 {
      margin: 0 0 8px;
      font-size: clamp(24px, 4vw, 38px);
    }

    .overlay p {
      margin: 0;
      color: var(--muted);
      line-height: 1.7;
    }

    @media (max-width: 760px) {
      body { padding: 14px; }
      .game-shell { grid-template-columns: 1fr; }
      .side-card { order: -1; }
      .mobile-pad { display: grid; }
      .tips { font-size: 14px; }
    }
  </style>
</head>
<body>
  <main class="game-shell" aria-label="贪吃蛇小游戏">
    <section class="board-card">
      <canvas id="board" width="560" height="560" aria-label="贪吃蛇游戏区域"></canvas>
      <div class="overlay" id="overlay">
        <div>
          <h2 id="overlayTitle">准备开始</h2>
          <p id="overlayText">按空格键或点击开始按钮，控制小蛇吃掉莓果。</p>
        </div>
      </div>
    </section>

    <aside class="side-card">
      <h1>贪吃蛇</h1>
      <p class="subtitle">经典玩法，现代霓虹风格。吃到食物会加分、变长并逐步提速。</p>

      <div class="stats">
        <div class="stat">
          <span class="label">当前分数</span>
          <span class="value" id="score">0</span>
        </div>
        <div class="stat">
          <span class="label">最高分</span>
          <span class="value" id="best">0</span>
        </div>
      </div>

      <div class="controls">
        <button id="startBtn">开始 / 暂停</button>
        <button class="secondary" id="resetBtn">重新开始</button>
      </div>

      <ul class="tips">
        <li>键盘：方向键或 WASD 控制移动。</li>
        <li>空格键：开始、暂停或继续游戏。</li>
        <li>手机：使用下方方向按钮控制。</li>
      </ul>

      <div class="mobile-pad" aria-label="手机方向控制">
        <button class="up" data-dir="up" aria-label="向上">▲</button>
        <button class="left" data-dir="left" aria-label="向左">◀</button>
        <button class="down" data-dir="down" aria-label="向下">▼</button>
        <button class="right" data-dir="right" aria-label="向右">▶</button>
      </div>
    </aside>
  </main>

  <script>
    const canvas = document.getElementById('board');
    const ctx = canvas.getContext('2d');
    const scoreEl = document.getElementById('score');
    const bestEl = document.getElementById('best');
    const overlay = document.getElementById('overlay');
    const overlayTitle = document.getElementById('overlayTitle');
    const overlayText = document.getElementById('overlayText');
    const startBtn = document.getElementById('startBtn');
    const resetBtn = document.getElementById('resetBtn');

    const gridSize = 20;
    const tileCount = canvas.width / gridSize;
    const initialSnake = [
      { x: 9, y: 12 },
      { x: 8, y: 12 },
      { x: 7, y: 12 }
    ];

    let snake;
    let food;
    let direction;
    let queuedDirection;
    let score;
    let best = Number(localStorage.getItem('snake-best-score') || 0);
    let status;
    let timer;
    let speed;

    function resetGame() {
      snake = initialSnake.map(part => ({ ...part }));
      direction = { x: 1, y: 0 };
      queuedDirection = direction;
      score = 0;
      speed = 130;
      status = 'ready';
      food = createFood();
      updateStats();
      showOverlay('准备开始', '按空格键或点击开始按钮，控制小蛇吃掉莓果。');
      draw();
      stopLoop();
    }

    function startGame() {
      if (status === 'gameover') {
        resetGame();
      }
      status = status === 'running' ? 'paused' : 'running';
      if (status === 'running') {
        hideOverlay();
        runLoop();
      } else {
        showOverlay('已暂停', '按空格键或点击开始按钮继续。');
        stopLoop();
      }
    }

    function runLoop() {
      stopLoop();
      timer = window.setInterval(tick, speed);
    }

    function stopLoop() {
      if (timer) {
        window.clearInterval(timer);
        timer = null;
      }
    }

    function tick() {
      direction = queuedDirection;
      const head = { x: snake[0].x + direction.x, y: snake[0].y + direction.y };

      if (isWallHit(head) || isSnakeHit(head)) {
        gameOver();
        return;
      }

      snake.unshift(head);
      if (head.x === food.x && head.y === food.y) {
        score += 10;
        if (score > best) {
          best = score;
          localStorage.setItem('snake-best-score', String(best));
        }
        speed = Math.max(62, speed - 4);
        food = createFood();
        runLoop();
      } else {
        snake.pop();
      }

      updateStats();
      draw();
    }

    function createFood() {
      let candidate;
      do {
        candidate = {
          x: Math.floor(Math.random() * tileCount),
          y: Math.floor(Math.random() * tileCount)
        };
      } while (snake.some(part => part.x === candidate.x && part.y === candidate.y));
      return candidate;
    }

    function setDirection(next) {
      const directions = {
        up: { x: 0, y: -1 },
        down: { x: 0, y: 1 },
        left: { x: -1, y: 0 },
        right: { x: 1, y: 0 }
      };
      const wanted = directions[next];
      if (!wanted) return;
      const isReverse = wanted.x + direction.x === 0 && wanted.y + direction.y === 0;
      if (!isReverse) queuedDirection = wanted;
      if (status === 'ready') startGame();
    }

    function isWallHit(point) {
      return point.x < 0 || point.y < 0 || point.x >= tileCount || point.y >= tileCount;
    }

    function isSnakeHit(point) {
      return snake.some(part => part.x === point.x && part.y === point.y);
    }

    function gameOver() {
      status = 'gameover';
      stopLoop();
      showOverlay('游戏结束', '最终得分：' + score + '。点击重新开始再来一局！');
      draw();
    }

    function updateStats() {
      scoreEl.textContent = score;
      bestEl.textContent = best;
    }

    function showOverlay(title, text) {
      overlayTitle.textContent = title;
      overlayText.textContent = text;
      overlay.hidden = false;
    }

    function hideOverlay() {
      overlay.hidden = true;
    }

    function draw() {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      drawGrid();
      drawFood();
      drawSnake();
    }

    function drawGrid() {
      ctx.fillStyle = '#020617';
      ctx.fillRect(0, 0, canvas.width, canvas.height);
      ctx.strokeStyle = 'rgba(148, 163, 184, 0.08)';
      ctx.lineWidth = 1;
      for (let line = 0; line <= canvas.width; line += gridSize) {
        ctx.beginPath();
        ctx.moveTo(line, 0);
        ctx.lineTo(line, canvas.height);
        ctx.stroke();
        ctx.beginPath();
        ctx.moveTo(0, line);
        ctx.lineTo(canvas.width, line);
        ctx.stroke();
      }
    }

    function drawSnake() {
      snake.forEach((part, index) => {
        const inset = index === 0 ? 2 : 3;
        ctx.fillStyle = index === 0 ? '#86efac' : '#22c55e';
        roundRect(part.x * gridSize + inset, part.y * gridSize + inset, gridSize - inset * 2, gridSize - inset * 2, 6);
        ctx.fill();
      });
    }

    function drawFood() {
      const centerX = food.x * gridSize + gridSize / 2;
      const centerY = food.y * gridSize + gridSize / 2;
      const gradient = ctx.createRadialGradient(centerX, centerY, 2, centerX, centerY, 11);
      gradient.addColorStop(0, '#fecdd3');
      gradient.addColorStop(1, '#fb7185');
      ctx.fillStyle = gradient;
      ctx.beginPath();
      ctx.arc(centerX, centerY, 8, 0, Math.PI * 2);
      ctx.fill();
    }

    function roundRect(x, y, width, height, radius) {
      ctx.beginPath();
      ctx.moveTo(x + radius, y);
      ctx.arcTo(x + width, y, x + width, y + height, radius);
      ctx.arcTo(x + width, y + height, x, y + height, radius);
      ctx.arcTo(x, y + height, x, y, radius);
      ctx.arcTo(x, y, x + width, y, radius);
      ctx.closePath();
    }

    startBtn.addEventListener('click', startGame);
    resetBtn.addEventListener('click', resetGame);

    document.querySelectorAll('[data-dir]').forEach(button => {
      button.addEventListener('click', () => setDirection(button.dataset.dir));
    });

    window.addEventListener('keydown', event => {
      const keyMap = {
        ArrowUp: 'up',
        KeyW: 'up',
        ArrowDown: 'down',
        KeyS: 'down',
        ArrowLeft: 'left',
        KeyA: 'left',
        ArrowRight: 'right',
        KeyD: 'right'
      };
      if (event.code === 'Space') {
        event.preventDefault();
        startGame();
        return;
      }
      const mapped = keyMap[event.code];
      if (mapped) {
        event.preventDefault();
        setDirection(mapped);
      }
    });

    resetGame();
  </script>
</body>
</html>`
