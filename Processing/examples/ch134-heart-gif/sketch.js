const heart = [];
const totalFrames = 240;
let counter = 0;

let angle = 0;
let c;

function setup() {
  createCanvas(windowWidth, windowHeight);
}

function draw() {
  draw2();
}

function draw1() {
  const percent = float(counter % totalFrames) / totalFrames;
  render(percent);
  counter++;
}

function render(percent) {
  background(0);
  translate(width / 2, height / 2);
  stroke(255);
  strokeWeight(4);
  fill(150, 0, 100);
  beginShape();
  for (let v of heart) {
    const a = map(percent, 0, 1, 0, TWO_PI * 2);
    const r = map(sin(a), -1, 1, height / 80, height / 40);
    vertex(r * v.x, r * v.y);
  }
  endShape();

  if (percent < 0.5) {
    const a = map(percent, 0, 0.5, 0, TWO_PI);
    const x = 16 * pow(sin(a), 3);
    const y = -(13 * cos(a) - 5 * cos(2 * a) - 2 * cos(3 * a) - cos(4 * a));
    heart.push(createVector(x, y));
  } else {
    heart.splice(0, 1);
  }
}

function draw2() {
  background(180, 0, 50, 30);
  //background(0);
  //save(c,`HEART-${counter}.jpg`);
  translate(width / 2, height / 2);
  fill(255, 0, 50);
  strokeWeight(4);
  noFill();
  stroke(255, 0, 50);
  let scl = map(sin(angle), 0, 1, 0.9, 1);
  scale(scl);

  if (angle < TWO_PI) {
    let r = map(pow(sin(angle), 1), 0, 1, 8, 10);
    let x = -1.5 * (16 * pow(sin(angle), 3));
    let y =
      -1.5 *
      (13 * cos(angle) -
        5 * cos(2 * angle) -
        2 * cos(3 * angle) -
        cos(4 * angle));
    heart.push(createVector(r * x, r * y));
  }
  if (angle > TWO_PI) {
    heart.shift();
    if (heart.length == 0) {
      angle = 0;
    }
  }

  beginShape();
  for (let pt of heart) {
    vertex(pt.x, pt.y);
  }
  endShape();

  angle += 0.05;
  counter += 1;
}
