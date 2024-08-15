// 箭塔模型
class Fort {
  constructor(x, y) {
    this.ox = x;
    this.oy = y;
    this.mx = 0;
    this.my = 0;
    this.angle = 0;
    this.color = "white";
    this.radiate = 0;
  }
  update() {
    this.mx = mouseX - this.ox;
    this.my = mouseY - this.oy;
    push();
    translate(this.ox, this.oy);
    if (this.my < 0) {
      this.angle = -atan(this.mx / this.my);
      rotate(this.angle);
    } else {
      this.angle = PI - atan(this.mx / this.my);
      rotate(this.angle);
    }
  }
  show() {
    if (this.radiate > 0 && this.radiate < 300) {
      stroke(r, g, 0);
      fill(r - 100, g - 100, 0, 150);
      ellipse(0, 0, this.radiate);
      this.radiate += 2;
      r -= 1.7;
      g -= 1.7;
    } else {
      this.radiate = 0;
      r = 255;
      g = 255;
    }
    stroke(this.color);
    noFill();
    ellipse(0, 0, ft_size);
    stroke(200);
    line(0, 0, 0, -bow_len);
    arc(0, 0, 1.5 * ft_size, 1.5 * ft_size, (-3 * PI) / 8, (11 * PI) / 8);
    triangle(
      0,
      -bow_len,
      -arw_w * tan(PI / 12),
      -arw_h,
      arw_w * tan(PI / 12),
      -arw_h
    );
    pop();
  }
}
