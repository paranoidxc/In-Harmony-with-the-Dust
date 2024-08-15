class Game {
  constructor(sx, sy) {
    this.sx = sx;
    this.sy = sy;
    this.score = 0;
    this.hp = 100;
    this.timer = 300;
    this.status = 0;
  }
  start() {
    noFill();
    stroke(255);
    rect(w / 2 - 50, h / 2 - 35, 100, 50);
    stroke("yellow");
    rect(w / 2 - 55, h / 2 - 40, 110, 60);
    textAlign(CENTER);
    textSize(24);
    noStroke();
    fill(255);
    text("START", w / 2, h / 2);
  }
  play() {
    noStroke();
    fill(255, 255, 255, 200);
    textStyle(BOLD);
    textSize(w / 32);
    textAlign(LEFT);
    text("Score：" + this.score, this.sx, this.sy);
    textAlign(RIGHT);
    text("HP：" + this.hp + "%", w - this.sx, this.sy);
  }
  end() {
    textAlign(CENTER);
    textSize(60);
    fill(255);
    stroke("red");
    strokeWeight(8);
    text("You Failed !", w / 2, h / 2);
    textSize(16);
    noStroke();
    strokeWeight(1);
    text("Click to restart.", w / 2, h / 2 + 30);
  }
  show() {
    switch (this.status) {
      case 0:
        this.start();
        break;
      case 1:
        this.play();
        break;
      case 2:
        this.end();
        break;
    }
  }
}
