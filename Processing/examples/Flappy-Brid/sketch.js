let bird;
let blocks = [];
let score;
let birdImg;
let bgImg;
let tubeImg;
let pipe;
let highestScore;
let pointSFX;
let jumpSFX;
let hitSFX;
let gameOn;
let gamePaused;
let myFont;
let curFrameRate;

function preload() {
  birdImg = loadImage("images/bat.png");
  bgImg = loadImage("images/batbg.png");
  pointSFX = loadSound("sfx/sfx_point.wav");
  jumpSFX = loadSound("sfx/sfx_jump.wav");
  hitSFX = loadSound("sfx/sfx_hit.wav");
  myFont = loadFont("./font/flappy.TTF");
}

function setup() {
  createCanvas(600, 400);
  angleMode(DEGREES);
  highestScore = 0;
  frameRate(30);
  background(255);
  textFont(myFont);
  push();
  fill(255);
  stroke(0);
  strokeWeight(3);
  textSize(50);
  textAlign(CENTER, CENTER);
  text("FLAPPY BIRD", width / 2, height / 3);
  textSize(20);
  text("press Space to play", width / 2, height / 1.5);
  pop();
  // console.log('bruh');
  noLoop();
  gameOn = false;
  gamePaused = false;
  curFrameRate = floor(frameRate());
}

//looping function
function draw() {
  if (gameOn && !gamePaused) {
    background(255);
    bird.update();
    bird.show();
    for (let block of blocks) {
      block.update();
      block.showTop();
      block.showBottom();
      if (block.passed) {
        //pointSFX.play();
        score++;
        block.passed = false;
        block.runCheck = false;
        // console.log("SCORE: "+score);
      }
      if (block.offScreen()) {
        blocks.shift();
        let block = new Blocks(1);
        blocks.push(block);
      }
      if (block.birdHit(bird)) {
        // hitSFX.play();
        gameOver();
        break;
      }
    }
    bird.hitBottom();
    bird.hitTop();
    showScore();

    push();
    fill(255);
    stroke(0);
    strokeWeight(3);
    textSize(30);
    text(curFrameRate, 20, 40);
    pop();
    if (frameCount % 10 == 0) {
      curFrameRate = floor(frameRate());
    }
  }
}

function keyPressed() {
  if (key === " " && gameOn && !gamePaused) {
    bird.click();
    //jumpSFX.play();
  }
  if (key === " ") {
    if (!gameOn) {
      // console.log("New Game");
      gameOn = true;
      gamePaused = false;
      restart();
      loop();
    }
  }

  if (key === "r") {
    restart();
    loop();
    gameOn = true;
    gamePaused = false;
  }

  if (key === "p") {
    if (!gamePaused && gameOn) {
      gamePaused = true;
      push();
      textSize(35);
      fill(255);
      stroke(0);
      strokeWeight(3);
      textAlign(CENTER, CENTER);
      text(
        `Paused

  Press P to resume
   or R to restart`,
        width / 2,
        175
      );
      pop();
      noLoop();
    } else if (gamePaused && gameOn) {
      loop();
      gamePaused = false;
    }
  }
}

function mouseClicked() {
  if (gameOn && !gamePaused) {
    bird.click();
    //jumpSFX.play();
  } else {
    restart();
    loop();
    gameOn = true;
  }
}

function restart() {
  bird = new Bird();
  frameRate(60);
  //if (frameCount % 100 == 0) {
  // blocks.push(new Blocks());
  //for (let i = 0; i < 5; i++) blocks[i] = new Blocks(i);
  //}
  for (let i = 0; i < 5; i++) blocks[i] = new Blocks(i);
  score = 0;
  gameOn = true;
}

function showScore() {
  push();
  fill(255);
  stroke(0);
  strokeWeight(3);
  textSize(30);
  text(score, width - 40, 40);
  pop();
}

function gameOver() {
  noLoop();
  push();
  fill(255);
  stroke(0);
  strokeWeight(3);
  textSize(50);
  textAlign(CENTER, CENTER);
  highestScore = max(highestScore, score);

  text(`Score : ${score}`, width / 2, height / 4);
  text(`Highest : ${highestScore}`, width / 2, height / 2.5);

  textSize(35);
  text("Press R to restart", width / 2, height / 1.5);
  pop();
  gameOn = false;
}
