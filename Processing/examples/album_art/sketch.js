let gap = 10;
let circleNum = 40;
let circleSize = 20;
let angle = 0;
let pointNum = 50;
let rectSize = 600;

function setup() {
  createCanvas(windowWidth, windowHeight);
  angleMode(RADIANS);
}

function draw() {
  background("black");
  noCursor();
  noStroke();
  fill("yellow");
  circle(mouseX, mouseY, 5);

  //main
  push();
  translate(width / 2, height / 2);
  rotate(angle);
  angle = angle + map(mouseX, 0, width, -0.1, 0.1);
  //text(angle, 25, 25);
  noFill();
  stroke("white");
  strokeWeight(1);
  for (let i = 0; i < circleNum; i++) {
    arc(
      0,
      0,
      circleSize + gap * i,
      circleSize + gap * i,
      angle * i,
      360 - angle / 2
    );
  }
  pop();

  let biteSize = PI / 16;
  let startAngle = biteSize * sin(frameCount * 0.1) + biteSize;
  let endAngle = TWO_PI - startAngle;

  // Draw the arc.
  arc(width / 2 - 84, height - 140, 16, 16, startAngle, endAngle, PIE);

  // title
  push();
  translate(width / 2, height - 140);
  textFont("Arial");
  textSize(15);
  textAlign(CENTER, CENTER);
  text("Paranoid.XiaoChuan", 0, 0);
  textSize(10);
  text("........ ... .....", 0, 20);
  pop();

  push();
  //border
  translate(width / 2, height / 2);
  noFill();
  stroke("white");
  strokeWeight(2);
  rectMode(CENTER);
  rect(0, 0, rectSize, rectSize);

  // random noise
  stroke("white");
  strokeWeight(1);
  for (let i = 0; i < pointNum; i++) {
    point(
      random(-rectSize / 2, rectSize / 2),
      random(-rectSize / 2, rectSize / 2)
    );
  }

  pop();
}
