class Bird {
  constructor() {
    this.x = 50;
    this.y = height / 5;
    this.r1 = 15;
    this.r2 = 15;
    this.gravity = 0.7;
    this.speed = 0;
    this.drag = 0.98;
    this.lift = -8;
    this.speedLimit = 20;
    this.angle = 0;
    this.clickAngle = 0;
  }

  show() {
    //absurd parameters to enhance symbol
    fill(0);
    push();
    translate(this.x, this.y);
    stroke(0);
    strokeWeight(2);
    fill(255);
    if (this.clickAngle > 0) {
      console.log("");
      this.angle += 20;
      this.clickAngle -= 1;
    } else {
      this.angle += 2;
    }
    rotate(this.angle);
    rect(-this.r1, -this.r2, this.r1 * 2, this.r2 * 2);
    pop();
  }

  update() {
    this.speed += this.gravity;
    this.y += this.speed;
    this.speed *= this.drag;
  }

  click() {
    this.speed = this.lift;
    this.clickAngle = 10;
  }

  hitBottom() {
    if (this.y >= height - this.r2) {
      this.y = height - this.r2;
      gameOver();
      //hitSFX.play();
      // this.speed = -5;
    }
  }

  hitTop() {
    if (this.y <= 0 + this.r2) {
      this.speed = 5;
      this.y = this.r2 + 10;
    }
  }
}
