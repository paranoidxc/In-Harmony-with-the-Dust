// 箭矢模型
class Arrow {
  constructor(bl, angle) {
    // 获取箭塔弓箭顶端坐标
    this.x0 = w / 2 - bl * sin(angle);
    this.y0 = h / 2 + bl * cos(angle);
    this.x = this.x0;
    this.y = this.y0;
    this.angle = angle;
    // 箭矢沿变换坐标后的y轴移动
    this.dy = 0;
    this.range = -200;
    this.timer = 60;
  }
  update() {
    push();
    translate(this.x0, this.y0);
    rotate(this.angle);
    if (
      isHit(this.x, this.y) == true ||
      this.dy <= this.range ||
      this.timer < 60
    ) {
      this.timer -= 1;
      this.dy -= 0;
    } else if (this.timer == 60) {
      this.dy -= 3;
      this.x = this.x0 - (this.dy - bow_len) * sin(this.angle);
      this.y = this.y0 + (this.dy - bow_len) * cos(this.angle);
    }
  }
  show() {
    if (this.timer < 60) {
      stroke(4 * this.timer, 4 * this.timer, 0);
      line(0, this.dy, 0, this.dy - bow_len);
      triangle(
        0,
        this.dy - bow_len,
        -arw_w * tan(PI / 12),
        this.dy - arw_h,
        arw_w * tan(PI / 12),
        this.dy - arw_h
      );
    } else {
      stroke(240, 240, 0);
      line(0, this.dy, 0, this.dy - bow_len);
      triangle(
        0,
        this.dy - bow_len,
        -arw_w * tan(PI / 12),
        this.dy - arw_h,
        arw_w * tan(PI / 12),
        this.dy - arw_h
      );
    }
    pop();
  }
}
