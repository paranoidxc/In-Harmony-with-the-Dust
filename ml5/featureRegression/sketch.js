/*
anything time you use a machine learning model
you want to ask youself the question

what data was used to train this model
who trained this model
what context and for what reason who trained this model
*/

let mobilenet;
let classifier;
let video;
let value = [];
let slider;
let addButton;
let trainButton;

let ukeButton;
let whistleButton;

function preload() {}

function modelReady() {
  console.log("Model is ready");
}

function videoReady() {
  console.log("Video is ready");
}

function whileTraining(loss) {
  if (loss == null) {
    console.log("Training Complete");
    classifier.predict(gotResult);
  } else {
    console.log(loss);
  }
}

function gotResult(error, result) {
  if (error) {
    console.error(error);
  } else {
    console.log(result);
    value = result;
    classifier.predict(gotResult);
  }
}

function setup() {
  value.value = 0;
  background(0);
  createCanvas(640, 560);
  video = createCapture(VIDEO);
  video.hide();
  mobilenet = ml5.featureExtractor("MobileNet", modelReady);
  classifier = mobilenet.regression(video, videoReady);

  slider = createSlider(0, 1, 0.5, 0.01);

  addButton = createButton("add Example image");
  addButton.mousePressed(function () {
    classifier.addImage(slider.value());
  });

  trainButton = createButton("train");
  trainButton.mousePressed(function () {
    classifier.train(whileTraining);
  });

  saveButton = createButton("save");
  saveButton.mousePressed(function () {
    classifier.save();
  });
}

function draw() {
  background(0);
  image(video, 0, 0);

  fill(225, 0, 200);
  rect(value.value * width, height / 2, 50, 50);

  fill(255);
  stroke(0);
  textSize(20);

  text(value.value, 10, height - 40);
}
