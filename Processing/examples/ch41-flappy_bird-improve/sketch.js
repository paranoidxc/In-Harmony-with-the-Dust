function preload() {
  //img = loadImage("")
}

var bird;
var pipes = [];
var mic;
var sliderTop;
var sliderBottom;
var clapping = false;

function setup() {
  createCanvas(400, 600);
  mic = new p5.AudioIn();
  mic.start();
  bird = new Bird();
  pipes.push(new Pipe());
  sliderTop = createSlider(0, 1, 0.2, 0.01);
  sliderBottom = createSlider(0, 1, 0.2, 0.01);
}

function draw() {
  background(0);
  bird.update();
  bird.show();

  var vol = mic.getLevel();

  if (frameCount % 100 == 0) {
    pipes.push(new Pipe());
  }

  fill(0, 255, 0);
  var y = map(vol, 0, 1, height, 0);
  rect(width - 50, y, 50, height - y);

  var thresholdT = sliderTop.value();
  var thresholdB = sliderBottom.value();
  console.log("vol", vol, "thresholdT", thresholdT, "thresholdB", thresholdB);
  if (vol > thresholdT && !clapping) {
    bird.up();
    clapping = true;
  }

  if (vol < thresholdB) {
    clapping = false;
  }

  var ty = map(thresholdT, 0, 1, height, 0);
  stroke(255, 0, 0);
  strokeWeight(4);
  line(width - 50, ty, width, ty);

  var by = map(thresholdB, 0, 1, height, 0);
  stroke(0, 0, 255);
  strokeWeight(4);
  line(width - 50, by, width, by);

  for (var i = pipes.length - 1; i >= 0; i--) {
    pipe = pipes[i];
    pipe.update();
    pipe.show();

    if (pipes[i].hits(bird)) {
    }

    if (pipe.offscreen()) {
      pipes.splice(i, 1);
    }
  }
}

function keyPressed() {
  if (key == " ") {
    bird.up();
  }
}
