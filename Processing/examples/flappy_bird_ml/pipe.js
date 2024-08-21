class Pipe {
  constructor() {
    this.spacing = 100;
    this.top = random(height - this.spacing);
    this.bottom = this.top + this.spacing;
    this.x = width;
    this.w = 20;
    this.velocity = 2;
  }

  collidesBottom(bird) {
    //圆的半径
    var radius = bird.r;
    //圆形中心与矩形中心的相对坐标
    var x = bird.x - this.x;
    var y = bird.y - this.top;

    var minX = Math.min(x, this.w / 2);
    var maxX = Math.max(minX, -this.w / 2);
    var minY = Math.min(y, this.top / 2);
    var maxY = Math.max(minY, -this.top / 2);

    if ((maxX - x) * (maxX - x) + (maxY - y) * (maxY - y) <= radius * radius) {
      return true;
    } else {
      return false;
    }
  }

  collidesTop(bird) {
    //圆的半径
    var radius = bird.r;
    //圆形中心与矩形中心的相对坐标
    var x = bird.x - this.x;
    var y = bird.y - this.bottom;

    var minX = Math.min(x, this.w / 2);
    var maxX = Math.max(minX, -this.w / 2);
    var minY = Math.min(y, (height - this.bottom) / 2);
    var maxY = Math.max(minY, -(height - this.bottom) / 2);

    if ((maxX - x) * (maxX - x) + (maxY - y) * (maxY - y) <= radius * radius) {
      return true;
    } else {
      return false;
    }
  }

  collides(bird) {
    /*
    if (this.collidesBottom(bird) || this.collidesTop(bird)) {
      return true;
    }
    return false;
    */
    return this.collidesOld(bird);
  }

  collidesOld(bird) {
    let birdCenterY = bird.y - bird.r;
    let birdCenterX = bird.x - bird.r;

    let verticalCollision = birdCenterY < this.top || birdCenterY > this.bottom;
    let horizontalCollision =
      birdCenterX > this.x && birdCenterX < this.x + this.w;
    if (verticalCollision && horizontalCollision) {
      console.log(
        "ver birdCenterY < this.top || birdCenterY > this.bottom",
        `ver ${birdCenterY} < ${this.top} || ${birdCenterY} > ${this.bottom}`
      );

      console.log(
        "hor  = birdCenterX > this.x && birdCenterX< this.x + this.w ",
        `hor  = ${birdCenterX} > ${this.x} && ${birdCenterX} < ${this.x} + ${this.w} `
      );
      return true;
    }
    return false;
  }

  show() {
    fill(0);
    noStroke();
    rect(this.x, 0, this.w, this.top);
    rect(this.x, this.bottom, this.w, height - this.bottom);
  }

  update() {
    this.x -= this.velocity;
  }

  offscreen() {
    return this.x < -this.w;
  }
}
