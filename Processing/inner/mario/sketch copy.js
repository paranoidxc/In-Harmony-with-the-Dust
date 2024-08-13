let gravity = 3;
let damping = 0.9;
let jumpSpeed = 85;
let moveSpeed = 0.22;
let marioright;
let windowWidth = 600;
let windowHeight = 720;
let startX = 100;
let startY = windowHeight - 78;
let myMario;

function preload() {
  marioright = loadImage("assets/marioright.png");
  marioleft = loadImage("assets/marioleft.png");
  mariowalkleft = loadImage("assets/mariowalkleft.png");
  mariowalkleft2 = loadImage("assets/mariowalkleft2.png");
  mariowalkright = loadImage("assets/mariowalkright.png");
  mariowalkright2 = loadImage("assets/mariowalkright2.png");
  // LtoR = loadImage("LtoR.png");
  // RtoL = loadImage("RtoL.png");
  // duckleft = loadImage("duckleft.png");
  // duckright = loadImage("duckright.png");
  // leftjump = loadImage("leftjump.png");
  // rightjump = loadImage("rightjump.png");
  // leftfall = loadImage("leftfall.png");
  // rightfall = loadImage("rightfall.png");
  // rightup = loadImage("rightup.png");
  // leftup = loadImage("leftup.png");
  // fastleft = loadImage("fastleft.png");
  // fastleft2 = loadImage("fastlCeft2.png");
  // fastleft3 = loadImage("fastleft3.png");
  // fastright = loadImage("fastright.png");
  // fastright2 = loadImage("fastright2.png");
  // fastright3 = loadImage("fastright3.png");
  // flyright = loadImage("flyright.png");
  // flyleft = loadImage("flyleft.png");
  // bg = loadImage("bg.png");
  // bg2 = loadImage("bg2.jpg");
  // ground = loadImage("ground.png");
  // ground2 = loadImage("ground2.png");
  // pipe = loadImage("pipe.png");
  // coin = loadImage("coin.png");
}

function setup() {
  createCanvas(windowWidth, windowHeight);
  this.hindernisX = 133;
  this.hindernisY = 233;
  myMario = new Mario(
    startX,
    startY,
    marioright,
    this.hindernisX,
    this.hindernisY
  );
}

function draw() {
  background("#6185f8");
  rect(this.hindernisX, this.hindernisY, 23, 56);

  fill("#954b0c");
  noStroke();
  rect(0, windowHeight - 50, windowWidth, windowHeight);

  myMario.move();
  myMario.display();
}

function keyPressed() {
  if (keyCode == LEFT_ARROW) {
    myMario.moveLeft = true;
  } else if (keyCode == RIGHT_ARROW) {
    myMario.moveRight = true;
    marioimage = loadImage("mariowalkright.png");
  } else if (keyCode == UP_ARROW) {
    myMario.moveUp = true;
  }
}

function keyReleased() {
  if (keyCode == LEFT_ARROW) {
    myMario.moveLeft = false;
  } else if (keyCode == RIGHT_ARROW) {
    myMario.moveRight = false;
  } else if (keyCode == UP_ARROW) {
    myMario.moveUp = false;
  }
}
