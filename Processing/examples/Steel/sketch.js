let font;
let vehicles = [];
let radius = 10;

function preload() {
  font = loadFont("DMSerifDisplay-Regular.ttf");
}

function setup() {
  createCanvas(1000, 300);
  background(51);
  /*
  textFont(font);
  textSize(192);
  fill(255);
  noStroke();
  text("train", 100, 200);
  */

  /*
  var points = font.textToPoints("train", 100, 200, 192);
  for (let i = 0; i < points.length; i++) {
    let pt = points[i];
    let vehicle = new Vehicle(pt.x, pt.y);
    vehicles.push(vehicle);
  }
    */

  let cols = floor(width / 10);
  let rows = floor(height / 10);

  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      let centerX = j * radius + radius / 2;
      let centerY = i * radius + radius / 2;

      let vehicle = new Vehicle(centerX, centerY);
      vehicles.push(vehicle);
    }
  }
}

function draw() {
  background(0);
  for (let i = 0; i < vehicles.length; i++) {
    let v = vehicles[i];
    v.behaviors();
    v.update();
    v.show();
  }
}

function mouseDragged() {}
