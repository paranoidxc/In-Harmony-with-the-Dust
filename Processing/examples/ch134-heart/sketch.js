const heart = [];
let a = 0;

function preload() {
  //img = loadImage("")
}

function setup() {
  createCanvas(640, 360);
}

function draw() {
  background(0);
  translate(width / 2, height / 2);
  stroke(255);
  strokeWeight(8);
  fill(150, 0, 100);
  beginShape();
  for (let v of heart) {
    vertex(v.x, v.y);
  }
  endShape();

  beginShape();
  fill("yellow");
  for (let v of heart) {
    //console.log("y", v.y);
    if (v.y > 0) {
      vertex(v.x, v.y);
    }
  }
  endShape();

  let r = 10;
  let x = r * 16 * pow(sin(a), 3);
  let y = -r * (13 * cos(a) - 5 * cos(2 * a) - 2 * cos(3 * a) - cos(4 * a));
  heart.push(createVector(x, y));
  if (a > TWO_PI) {
    noLoop();
  }
  a += 0.1;
}

function drawSingle() {
  background(0);
  push();

  translate(width / 2, height / 2);

  fill("yellow");
  rect(-200, 50, 400, 100);

  noFill();
  stroke(255);
  beginShape();
  for (let a = 0; a < TWO_PI; a += 0.01) {
    let r = 10;
    let x = r * 16 * pow(sin(a), 3);
    let y = -r * (13 * cos(a) - 5 * cos(2 * a) - 2 * cos(3 * a) - cos(4 * a));
    vertex(x, y);
  }
  endShape(CLOSE);

  pop();
}
