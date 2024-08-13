/*
anything time you use a machine learning model
you want to ask youself the question

what data was used to train this model
who trained this model
what context and for what reason who trained this model
*/

let features;
let video;
let knn;
let labelP;
let ready = false;
let label = "";
let x, y;

function modelReady() {
  console.log("MobileNet is ready");
  knn = ml5.KNNClassifier();
  knn.load("model.json", function () {
    console.log("KNN data is ready");
    goClassify();
  });
}

function videoReady() {
  console.log("Video is ready");
}

function setup() {
  background(0);
  createCanvas(320, 240);
  video = createCapture(VIDEO);
  video.size(320, 240);
  video.style("transform", "scale(-1, 1)");
  //video.hide();
  features = ml5.featureExtractor("MobileNet", modelReady);
  labelP = createP("need training data");
  x = width / 2;
  y = height / 2;
}

function goClassify() {
  const logits = features.infer(video);
  knn.classify(logits, function (error, result) {
    if (error) {
      console.error(error);
    } else {
      //console.log(result);
      label = result.label;
      labelP.html(label);
      goClassify();
    }
  });
}

function keyPressed() {
  const logits = features.infer(video);
  if (key == "l") {
    knn.addExample(logits, "left");
    console.log("left");
  } else if (key == "r") {
    knn.addExample(logits, "right");
    console.log("right");
  } else if (key == "u") {
    knn.addExample(logits, "up");
    console.log("up");
  } else if (key == "d") {
    knn.addExample(logits, "down");
    console.log("down");
  } else if (key == "s") {
    knn.save("model.json");
  }
  //console.log(logits.dataSync());
}

function draw() {
  background(0);
  fill(255);
  ellipse(x, y, 36);
  if (label == "up") {
    y--;
  } else if (label == "down") {
    y++;
  } else if (label == "left") {
    x--;
  } else if (label == "right") {
    x++;
  }
  x = constrain(x, 0, width);
  y = constrain(y, 0, height);
  //image(video, 0, 0);
}
