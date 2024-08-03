let bullets = [];
let enemies = [];
let score = 0;
let life = 3;
let isHit = 0;

function setup() {
  createCanvas(400, 600);

  // spawn enemies
  for (let i = 0; i < 20; i++) {
    let enemy = newEnemy();
    enemies.push(enemy);
  }
}

function newEnemy() {
  return {
    x: random(0, width),
    y: random(-height, -height / 300),
  };
}

function draw() {
  background(51);
  if (isHit) {
    fill("red");
    isHit--;
    if (isHit <= 0) {
      isHit = 0;
    }
  }
  circle(mouseX, height - 20, 25);

  fill("white");
  for (let bullet of bullets) {
    bullet.y -= 10;
    circle(bullet.x, bullet.y, 10);
  }

  for (let enemy of enemies) {
    enemy.y += 2;
    rect(enemy.x, enemy.y, 10);
    if (enemy.y > height) {
      enemy.y = 0;
    }
  }

  for (let enemy of enemies) {
    for (let bullet of bullets) {
      if (dist(enemy.x, enemy.y, bullet.x, bullet.y) < 10) {
        enemies.splice(enemies.indexOf(enemy), 1);
        enemies.push(newEnemy());
        bullets.splice(bullets.indexOf(bullet), 1);
        score++;
        continue;
      }
    }

    if (dist(enemy.x, enemy.y, mouseX - 12.5, height - 20) < 17.5) {
      if (life > 0) {
        enemies.splice(enemies.indexOf(enemy), 1);
        life--;
        isHit = 5;
        enemies.push(newEnemy());
        continue;
      } else {
        text("You Died", width / 2 - 30, height / 2);
        noLoop();
      }
    }
  }

  textSize(20);
  fill("yellow");
  text("score: " + score, 10, 30);

  textSize(20);
  fill("red");
  for (i = 0; i < life; i++) {
    text("❤️", width - 24 - i * 24, 30);
  }
  fill("white");
}

function mousePressed() {
  let bullet = {
    x: mouseX,
    y: height - 50,
  };

  bullets.push(bullet);
}
