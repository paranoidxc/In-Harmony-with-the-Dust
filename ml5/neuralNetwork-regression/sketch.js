let model;
let targetLabel = "C";
let trainingData = [];

let state = "collection";
let wave;

let notes = {
  C: 261.6256,
  D: 293.6648,
  E: 329.6276,
  F: 349.2282,
  G: 391.9954,
  A: 440.0,
  B: 493.8833,
};

function setup() {
  createCanvas(400, 480);
  env = new p5.Envelope();
  env.setADSR(0.05, 0.1, 0.5, 1);
  env.setRange(1.2, 0);

  wave = new p5.Oscillator();
  wave.setType("sine");
  wave.start();
  wave.freq(440);
  wave.amp(env);

  let options = {
    inputs: ["x", "y"],
    outputs: ["frequency"],
    task: "regression",
    debug: "true",
  };
  background(200);
  model = ml5.neuralNetwork(options);
  //model.loadData("./cdefgab-notes.json", dataLoaded);
  /*
  const modelInfo = {
    model: "model/model.json",
    metadata: "model/model_meta.json",
    weights: "model/model.weights.bin",
  };
  model.load(modelInfo, modelLoaded);
  */
}

function modelLoaded() {
  console.log("modelLoaded");
}

function dataLoaded() {
  console.log("dataLoaded");
  let data = model.neuralNetworkData.data.raw;
  for (let i = 0; i < data.length; i++) {
    let inputs = data[i].xs;
    let target = data[i].ys;

    stroke(0);
    noFill();
    ellipse(inputs.x, inputs.y, 24);
    fill(0);
    noStroke();
    textAlign(CENTER, CENTER);
    text(target.label, inputs.x, inputs.y);
  }
  model.normalizeData();
  let options = {
    epochs: 200,
  };
  model.train(options, whileTraining, finishedTraining);
}

function keyPressed() {
  console.info("key", key);
  if (key == "t") {
    state = "trainning";
    model.normalizeData();
    let options = {
      epochs: 100,
    };
    model.train(options, whileTraining, finishedTraining);
  } else if (key == "s") {
    model.saveData("cdefgab-notes");
  } else if (key == "m") {
    model.save("cdefgab-notes");
  } else {
    targetLabel = key.toUpperCase();
  }
}

function whileTraining(epochs, loss) {
  console.log("epochs", epochs);
}

function finishedTraining() {
  console.log("finished training");
  state = "prediction";
}

function mousePressed() {
  let inputs = {
    x: mouseX,
    y: mouseY,
  };
  if (state == "collection") {
    let targetFrequency = notes[targetLabel];
    let target = {
      frequency: targetFrequency,
    };
    model.addData(inputs, target);

    stroke(0);
    noFill();
    ellipse(mouseX, mouseY, 24);
    fill(0);
    noStroke();
    textAlign(CENTER, CENTER);
    text(targetLabel, mouseX, mouseY);

    wave.freq(targetFrequency);
    env.play();
  } else if (state == "prediction") {
    model.predict(inputs, gotResult);
  }
}

function gotResult(error, result) {
  if (error) {
    console.log(error);
    return;
  }
  console.log(result);
  stroke(0);
  fill(0, 0, 255, 100);
  ellipse(mouseX, mouseY, 24);
  fill(0);
  noStroke();
  textAlign(CENTER, CENTER);
  text(floor(result[0].value), mouseX, mouseY);
  wave.freq(result[0].value);
  env.play();
}

function draw() {}
