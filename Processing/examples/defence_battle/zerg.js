// 虫子模型
class Zerg {
  constructor(x, y, size) {
    this.x0 = x;
    this.y0 = y;
    this.x = this.x0;
    this.y = this.y0;
    // 虫子沿变换坐标后的x轴移动
    this.dx = 0;
    this.size = size;
    this.tx = w / 2 - x;
    this.ty = h / 2 - y;
    this.angle =
      this.tx > 0 ? atan(this.ty / this.tx) : PI + atan(this.ty / this.tx);
    // 虫子与炮塔中心的距离
    this.distance = dist(x, y, this.tx, this.ty);
    this.timer = 60;
    this.alive = true;
    this.burst = false;
    this.destroy = false;
  }
  update() {
    translate(this.x0, this.y0);
    rotate(this.angle);
    this.distance = dist(this.dx, 0, this.tx / cos(this.angle), 0);
    if (this.distance <= 2 * this.size && this.burst == false) {
      this.burst = true;
    } else if (floor(this.distance) <= fort.radiate / 2) {
      this.destroy = true;
    }
    if (this.burst == true || this.alive == false || this.destroy == true) {
      this.dx += 0;
    } else {
      this.dx += random(1);
      this.x = w / 2 - this.distance * cos(this.angle);
      this.y = h / 2 - this.distance * sin(this.angle);
    }
  }
  dead() {
    if (this.timer > 0) {
      stroke(2 * this.timer, 2 * this.timer, 2 * this.timer);
      rect(this.dx - this.size / 2, -this.size / 2, this.size);
      this.timer -= 1;
      if (this.timer == 0) {
        info.score += 10;
      }
    }
  }
  boom() {
    if (this.timer > 0) {
      this.timer < 6 ? stroke("yellow") : stroke("red");
      if (frameCount % 3 == 0) {
        rect(this.dx - this.size / 2, -this.size / 2, this.size);
      }
      this.timer -= 1;
      if (this.timer == 0) {
        boomS.play();
        if (this.destroy == true) {
          info.score += 10;
        }
        if (this.burst == true) {
          info.hp -= 10;
          if (info.hp <= 0) {
            info.status = 2;
          }
        }
      }
    }
  }
  show() {
    noFill();
    if (this.alive === false) {
      this.dead();
    } else if (this.burst === true || this.destroy == true) {
      this.boom();
    } else {
      stroke("green");
      // 使虫子的中心点保持在x轴上
      rect(this.dx - this.size / 2, -this.size / 2, this.size);
    }
  }
}
