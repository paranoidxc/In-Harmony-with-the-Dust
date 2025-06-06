/*
anything time you use a machine learning model
you want to ask youself the question

what data was used to train this model
who trained this model
what context and for what reason who trained this model
*/

let vedio;
let detector;
let detections = [];
function preload() {
  detector = ml5.objectDetector("cocossd");
}

function modelLoaded() {
  console.log("Model Loaded!");
}

function gotDetections(error, results) {
  if (error) {
    console.error(error);
  }
  detections = results;
  detector.detect(video, gotDetections);
}

function setup() {
  createCanvas(640, 480);
  video = createCapture(VIDEO);
  video.hide();
  //video.size(640, 480);
  detector.detect(video, gotDetections);
}

function draw() {
  image(video, 0, 0);
  for (let i = 0; i < detections.length; i++) {
    let object = detections[i];
    stroke(0, 255, 0);
    strokeWeight(4);
    noFill();
    rect(object.x, object.y, object.width, object.height);
    noStroke();
    fill(0);
    textSize(24);
    text(object.label, object.x + 10, object.y + 24);
  }
}
