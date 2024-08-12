const density = "Ñ@#W$9876543210?!abc;:+=-,._          ";
let gloria;

function preload() {
  gloria = loadImage("gloria48.jpg");
}

function setup() {
  //createCanvas(400, 400);
  //noCanvas();
  //}

  //function draw() {
  background(0);
  //image(gloria, 0, 0, width, height);

  w = width / gloria.width;
  h = height / gloria.height;
  //console.log(width, gloria.width);
  gloria.loadPixels();

  const len = density.length;
  for (let i = 0; i < gloria.width; i++) {
    for (let j = 0; j < gloria.height; j++) {
      const pixelIndex = (i + j * gloria.width) * 4;
      const r = gloria.pixels[pixelIndex + 0];
      const g = gloria.pixels[pixelIndex + 1];
      const b = gloria.pixels[pixelIndex + 2];
      const a = gloria.pixels[pixelIndex + 3];

      const avg = (r + g + b) / 3;
      noStroke();
      //fill(avg);
      //square(i * w, j * h, w);
      fill(255);

      const charIndex = floor(map(avg, 0, 255, len, 0));
      textSize(w);
      text(density[charIndex], i * w, j * h);
    }
  }
}
