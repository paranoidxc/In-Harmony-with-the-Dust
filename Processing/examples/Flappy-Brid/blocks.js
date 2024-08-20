class Blocks {
  constructor(i) {
    this.gapSize = 130;
    //this.gapSize = 200;
    this.gapLocation = random(20 + this.gapSize, height - (20 + this.gapSize));
    this.top = this.gapLocation - this.gapSize / 2;
    this.bottom = this.gapLocation + this.gapSize / 2;
    this.x = width + 180 * i;
    this.speed = -2;
    this.w = 40;
    this.passed = false;
    this.runCheck = true;
    this.horizontalDivs = 5;
    this.windowSize = this.w / this.horizontalDivs;
    this.windowGap = 12;
    this.dist = this.windowSize * 2 + this.windowGap;
    this.windowsTop = floor(this.top / this.dist);
    this.windowsBottom = floor((height - this.bottom) / this.dist);
    this.colour = random(30, 60);
    this.max = floor((this.top - this.windowSize) / (this.windowSize * 3));
  }

  showTop() {
    strokeWeight(1);
    stroke(0);
    fill(255);
    rect(this.x, 0 - 1, this.w, this.top);

    for (let i = this.max; i > 0; i--) {
      fill(111);
      rect(
        this.x + this.windowSize,
        this.top - this.windowSize * 3 * i - this.windowSize,
        this.windowSize,
        this.windowSize * 2
      );
      rect(
        this.x + this.w - 2 * this.windowSize,
        this.top - this.windowSize * 3 * i - this.windowSize,
        this.windowSize,
        this.windowSize * 2
      );
    }
  }

  showBottom() {
    stroke(0);
    strokeWeight(1);
    fill(255);
    rect(this.x, this.bottom, this.w, height - this.bottom + 1);
    for (let i = 0; i < this.windowsBottom; i++) {
      fill(111);
      rect(
        this.x + this.windowSize,
        this.bottom +
          this.windowGap +
          i * (this.windowGap + this.windowSize * 2),
        this.windowSize,
        this.windowSize * 2
      );
      rect(
        this.x + this.windowSize * ((this.horizontalDivs + 1) / 2),
        this.bottom +
          this.windowGap +
          i * (this.windowGap + this.windowSize * 2),
        this.windowSize,
        this.windowSize * 2
      );
    }
  }

  update() {
    this.x += this.speed;
    if (this.runCheck) this.passed = bird.x - bird.r1 > this.x - this.w;
  }

  birdHit(bird) {
    if (bird.x + bird.r1 >= this.x && bird.x - bird.r1 <= this.x + this.w) {
      if (bird.y - bird.r2 <= this.top || bird.y + bird.r2 >= this.bottom)
        return true;
    }
    return false;
  }

  offScreen() {
    return this.x + this.w < 0;
  }
}
