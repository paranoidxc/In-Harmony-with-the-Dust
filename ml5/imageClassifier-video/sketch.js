let mobilenet;
let video;
let label = "";
let confidence = "";

/*
anything time you use a machine learning model
you want to ask youself the question

what data was used to train this model
who trained this model
what context and for what reason who trained this model
*/

function preload() {}

function modelReady() {
  console.log("Model is ready");
  mobilenet.classifyStart(video, gotResult);
}

function setup() {
  background(0);
  createCanvas(640, 560);
  video = createCapture(VIDEO);
  video.hide();
  mobilenet = ml5.imageClassifier("MobileNet", modelReady);
}

function draw() {
  background(0);
  image(video, 0, 0);

  fill(255);
  stroke(0);
  textSize(20);
  text(label, 10, height - 40);
  text(confidence, 10, height - 20);
}

function gotResult(results) {
  //console.log(results);

  label = "Label: " + results[0].label;
  confidence = "Confidence: " + nf(results[0].confidence, 0, 2);
}
