//PLAYER CLASS
class Mario {
  constructor(xin, yin, image_in, hinx, hiny) {
    this.hindernisX = hinx;
    this.hindernisY = hiny;
    this.xPo = xin;
    this.yPo = yin;
    this.xVel = 0;
    this.yVel = 0;
    this.speed = 2;
    this.moveLeft = false;
    this.moveRight = false;
    this.moveUp = false;
    this.onGround = false;
    this.marioimage = image_in;
  }

  display() {
    fill(255, 0, 0);
    noStroke();
    image(this.marioimage, this.xPo, this.yPo);
  }

  move() {
    if (this.moveLeft) {
      this.xVel = this.xVel - moveSpeed;
    }
    if (this.moveRight) {
      this.xVel = this.xVel + moveSpeed;
    }
    if (this.moveUp && this.onGround) {
      this.yVel = -jumpSpeed;
      this.onGround = false;
    }
    print(this.xPo);
    print(this.yPo);

    this.xPo = this.xPo + this.xVel;
    this.yPo = this.yPo + this.yVel;
    this.xVel = this.xVel * damping;
    this.yVel = this.yVel * damping;

    if (!this.onGround) {
      this.yVel = this.yVel + gravity;
    }

    if (this.yPo >= startY) {
      this.onGround = true;
      this.yPo = startY;
    }
    if (this.yPo <= this.hindernisY) {
      this.yPo = this.hindernisY;
      this.yVel = 0;
    }
    if (this.moveRight) {
      this.marioimage = mariowalkright;
    }

    if (this.moveLeft) {
      this.marioimage = mariowalkleft;
    } else if (this.moveRight == false && this.moveLeft == false) {
      this.marioimage = marioright;
    }
  }
}
