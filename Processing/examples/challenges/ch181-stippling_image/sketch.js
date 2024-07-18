let organic = true; // randomly moves points around

let points = [];
let delaunay, voronoi;
let started = false;

let gloria;

/** This function loads resources that will be used later. */
function preload() {
  gloria = loadImage("2.jpeg");
}

function setup() {
  createCanvas(800, 600);
  for (let i = 0; i < 6000; i++) {
    let x = random(width);
    let y = random(height);
    let col = gloria.get(x, y);
    //if (160random(100) > brightness(col)) {
    if (random(255) > brightness(col)) {
      points.push(createVector(x, y));
    } else {
      i--;
    }
  }
  delaunay = calculateDelaunay(points);
  voronoi = delaunay.voronoi([0, 0, width, height]);

  started = true;
}

function draw() {
  background(255);

  let polygons = voronoi.cellPolygons();
  let cells = Array.from(polygons);

  let centroids = new Array(cells.length);
  let weights = new Array(cells.length).fill(0);
  let weightsR = new Array(cells.length).fill(0);
  let weightsG = new Array(cells.length).fill(0);
  let weightsB = new Array(cells.length).fill(0);
  let counts = new Array(cells.length).fill(0);
  let avgWeights = new Array(cells.length).fill(0);
  for (let i = 0; i < centroids.length; i++) {
    centroids[i] = createVector(0, 0);
  }

  gloria.loadPixels();
  let delaunayIndex = 0;
  for (let i = 0; i < width; i++) {
    for (let j = 0; j < height; j++) {
      let index = (i + j * width) * 4;
      let r = gloria.pixels[index + 0];
      let g = gloria.pixels[index + 1];
      let b = gloria.pixels[index + 2];
      let bright = (r + g + b) / 3;
      let weight = 1 - bright / 255;
      delaunayIndex = delaunay.find(i, j, delaunayIndex);
      centroids[delaunayIndex].x += i * weight;
      centroids[delaunayIndex].y += j * weight;
      weights[delaunayIndex] += weight;
      weightsR[delaunayIndex] += r * r;
      weightsG[delaunayIndex] += g * g;
      weightsB[delaunayIndex] += b * b;
      counts[delaunayIndex]++;
    }
  }

  let maxWeight = 0;
  for (let i = 0; i < centroids.length; i++) {
    if (weights[i] > 0) {
      centroids[i].div(weights[i]);
      avgWeights[i] = weights[i] / (counts[i] || 1);
      if (avgWeights[i] > maxWeight) {
        maxWeight = avgWeights[i];
      }
    } else {
      centroids[i] = points[i].copy();
    }
  }

  for (let i = 0; i < points.length; i++) {
    if (organic && random(1) > 0.98)
      points[i].set(random(width), random(height));
    else points[i].lerp(centroids[i], 1);
  }

  for (let i = 0; i < points.length; i++) {
    let v = points[i];
    let index = (floor(v.x) + floor(v.y) * width) * 4;
    let r = gloria.pixels[index + 0];
    let g = gloria.pixels[index + 1];
    let b = gloria.pixels[index + 2];
    let sw = map(avgWeights[i], 0, maxWeight, 0, 12, true);
    let col = color(r, g, b);
    if (weights[i] > 0) {
      let ra = Math.sqrt(weightsR[i] / weights[i]);
      let ga = Math.sqrt(weightsG[i] / weights[i]);
      let ba = Math.sqrt(weightsB[i] / weights[i]);
      col = color(ra, ga, ba);
    }
    strokeWeight(0.5);
    stroke(0, 0, 0);
    fill(col);
    let poly = cells[i];
    beginShape();
    for (let i = 0; i < poly.length; i++) {
      vertex(poly[i][0], poly[i][1]);
    }
    endShape();
    point(v.x, v.y);
  }

  delaunay = calculateDelaunay(points);
  voronoi = delaunay.voronoi([0, 0, width, height]);
}

function calculateDelaunay(points) {
  let pointsArray = [];
  for (let v of points) {
    pointsArray.push(v.x, v.y);
  }
  return new d3.Delaunay(pointsArray);
}
