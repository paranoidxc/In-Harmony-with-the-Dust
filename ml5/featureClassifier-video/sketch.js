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
let label = [];
let confidence = "";

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
    classifier.classify(gotResult);
  } else {
    console.log(loss);
  }
}

function gotResult(error, results) {
  if (error) {
    console.error(error);
  } else {
    console.log(results);
    label = results;
    classifier.classify(gotResult);
  }
}

function setup() {
  background(0);
  createCanvas(640, 560);
  video = createCapture(VIDEO);
  video.hide();
  mobilenet = ml5.featureExtractor("MobileNet", modelReady);
  //classifier = mobilenet.classification(video, videoReady);
  classifier = mobilenet.regression(video, videoReady);

  ukeButton = createButton("ukulele");
  ukeButton.mousePressed(function () {
    classifier.addImage("ukulele");
  });

  whistleButton = createButton("whistle");
  whistleButton.mousePressed(function () {
    classifier.addImage("whistle");
  });

  trainButton = createButton("train");
  trainButton.mousePressed(function () {
    classifier.train(whileTraining);
  });
}

function draw() {
  background(0);
  image(video, 0, 0);

  fill(255);
  stroke(0);
  textSize(20);

  let txt = "test";
  if (label.length) {
    txt = "Label: " + label[0].label;
  }

  text(txt, 10, height - 40);
  text(confidence, 10, height - 20);
}
