let angle = 0;

function setup() {
  createCanvas(400, 400);
}

function drawDebug(x, y, radius) {
  drawingContext.setLineDash([5, 3]);
  noFill();
  stroke("grey");
  circle(x, y, 2 * radius);

  fill("grey");
  circle(x, y, 4);

  line(x, y, x, y - radius);

  drawingContext.setLineDash([]);
}

function rotateText(x, y, radius, txt) {
  // Comment the following line to hide debug objects
  drawDebug(x, y, radius);

  // Split the chars so they can be printed one by one
  chars = txt.split("");

  // Decide an angle
  charSpacingAngleDeg = 8;

  // https://p5js.org/reference/#/p5/textAlign
  textAlign(CENTER, BASELINE);
  textSize(15);
  fill("black");

  // https://p5js.org/reference/#/p5/push
  // Save the current translation matrix so it can be reset
  // before the end of the function
  push();

  // Let's first move to the center of the circle
  translate(x, y);

  // First rotate half back so that middle char will come in the center
  console.log(radians((-chars.length * charSpacingAngleDeg) / 2));
  angle = angle + map(mouseX, 0, width, -0.1, 0.1);
  let radiansx = radians((-chars.length * charSpacingAngleDeg) / 2) + angle;

  rotate(radiansx);
  //rotate(radians((-chars.length * charSpacingAngleDeg) / 2));

  for (let i = 0; i < chars.length; i++) {
    text(chars[i], 0, -radius);

    // Then keep rotating forward per character
    rotate(radians(charSpacingAngleDeg));
  }

  // Reset all translations we did since the last push() call
  // so anything we draw after this isn't affected
  pop();
}

function draw() {
  background(220);

  textToRotate = "草班台子";
  rotateText(200, 200, 150, textToRotate);

  textToRotate = "草班台子22222";
  rotateText(200, 200, 200, textToRotate);
}
