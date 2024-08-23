var x = 0;
var speed = 3;

function setup() {
  createCanvas(600, 400);
}

function draw() {
  bakgrunn = createGraphics(width, height);
  bakgrunn.background(0);
  bakgrunn.fill("red");
  bakgrunn.textSize(width / 3);
  bakgrunn.textStyle(BOLD);
  bakgrunn.textAlign(CENTER, CENTER);
  let t = nf(hour(), 2);
  let m = nf(minute(), 2);
  bakgrunn.text(t + ":" + m, width / 2, height * 0.5);

  //let farge = bakgrunn.get(300, 200);
  //console.log(farge);

  image(bakgrunn, 0, 0);
}
