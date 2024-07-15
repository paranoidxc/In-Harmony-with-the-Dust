function Drop() {
  this.x = random(width);
  this.y = random(-500, -50);
  this.z = random(0, 20);
  this.color = color(random(255), random(255), random(255));
  this.len = map(this.z, 0, 20, 10, 20);
  this.yspeed = map(this.z, 0, 20, 1, 20);

  this.fall = () => {
    this.y = this.y + this.yspeed;
    var grav = map(this.z, 0, 20, 0, 0.2);
    this.yspeed = this.yspeed + grav;

    if (this.y > height) {
      ellipse(this.x, this.y, 12, 12);
      this.y = random(-200, -100);
      this.yspeed = map(this.z, 0, 20, 4, 10);
    }
  };

  this.show = () => {
    var thick = map(this.z, 0, 20, 1, 3);
    strokeWeight(thick);
    fill(this.color);
    stroke(this.color);
    line(this.x, this.y, this.x, this.y + this.len);
  };
}
