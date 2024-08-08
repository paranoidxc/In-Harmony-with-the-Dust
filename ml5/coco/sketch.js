/*
anything time you use a machine learning model
you want to ask youself the question

what data was used to train this model
who trained this model
what context and for what reason who trained this model
*/

let img;
let detector;

function preload() {
  img = loadImage("images/1111.jpeg");
  detector = ml5.objectDetector("cocossd");
}

function modelLoaded() {
  console.log("Model Loaded!");
}

function gotDetections(error, results) {
  if (error) {
    console.error(error);
  }
  console.log(results);

  for (let i = 0; i < results.length; i++) {
    let object = results[i];
    stroke(0, 255, 0);
    strokeWeight(4);
    noFill();
    rect(object.x, object.y, object.width, object.height);
    noStroke();
    fill(255);
    textSize(24);
    text(object.label, object.x + 10, object.y + 24);
  }
}

function setup() {
  background(0);
  createCanvas(640, 480);
  image(img, 0, 0);
  detector.detect(img, gotDetections);
}

function draw() {}
