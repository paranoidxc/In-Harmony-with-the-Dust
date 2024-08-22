let shapeA;
let shapeB;

let isCol = false;

function setup() {
  createCanvas(640, 360);
  shapeA = {
    x: 320,
    y: 180,
    w: 100,
    h: 100,
    rotation: 0,
  };

  shapeB = {
    x: 100,
    y: 100,
    w: 40,
    h: 40,
    //rotation: 30,
    rotation: 0,
  };
}

function draw() {
  background(200);
  let ret = check();
  if (ret) {
    isCol = true;
  } else {
    isCol = false;
  }

  push();
  noFill();
  if (isCol) {
    stroke("red");
  } else {
    stroke("blue");
  }
  translate(shapeA.x, shapeA.y);
  rect(-shapeA.w / 2, -shapeA.h / 2, shapeA.w, shapeA.h);
  pop();

  push();
  noFill();
  if (isCol) {
    stroke("red");
  } else {
    stroke("purple");
  }
  translate(shapeB.x, shapeB.y);
  angleMode(DEGREES);
  rotate(shapeB.rotation);
  rect(-shapeB.w / 2, -shapeB.h / 2, shapeB.w, shapeB.h);
  pop();

  //angleMode(RADIANS);
}

function check() {
  let shortestDist = Number.MAX_VALUE;

  let verts2 = getTransformedVerts(shapeB);
  //console.log("vertsB", verts2);

  let verts1 = getTransformedVerts(shapeA);
  //console.log("vertsA", verts1);

  let vOffset = new SATPoint(shapeA.x - shapeB.x, shapeA.y - shapeB.y);

  //console.log("shapeB:x", shapeA.x + shapeA.w / 2);
  //console.log("shapeB:y", shapeA.y + shapeA.h / 2);

  // set up the result object
  let result = new SATCollisionInfo();
  result.shapeA = shapeA;
  result.shapeB = shapeB;
  result.shapeAContained = true;
  result.shapeBContained = true;

  for (let i = 0; i < verts1.length; i++) {
    let axis = _getPerpendicularAxis(verts1, i);
    //console.log("axis", axis);

    let polyARange = _projectVertsForMinMax(axis, verts1);
    //console.log("polyARange", polyARange);

    let polyBRange = _projectVertsForMinMax(axis, verts2);
    //console.log("polyBRange", polyBRange);

    var scalerOffset = _vectorDotProduct(axis, vOffset);
    //console.log("scalerOffset =", scalerOffset);
    polyARange.min += scalerOffset;
    polyARange.max += scalerOffset;

    // now check for a gap betwen the relative min's and max's
    if (
      polyARange.min - polyBRange.max > 0 ||
      polyBRange.min - polyARange.max > 0
    ) {
      // there is a gap - bail
      //console.log("GAP FOUND");
      return null;
    } else {
      //console.log("NOT FOUND");
    }

    let flipResultPositions = false;
    _checkRangesForContainment(
      polyARange,
      polyBRange,
      result,
      flipResultPositions
    );

    // calc the separation and store if this is the shortest
    let distMin = (polyBRange.max - polyARange.min) * -1;
    if (flipResultPositions) distMin *= -1;

    // check if this is the shortest by using the absolute val
    let distMinAbs = Math.abs(distMin);
    if (distMinAbs < shortestDist) {
      shortestDist = distMinAbs;

      result.distance = distMin;
      result.vector = axis;
    }
  }

  result.separation = new SATPoint(
    result.vector.x * result.distance,
    result.vector.y * result.distance
  );

  return result;
}

function _projectVertsForMinMax(axis, verts) {
  // note that we project the first point to both min and max
  let min = _vectorDotProduct(axis, verts[0]);
  let max = min;

  // now we loop over the remiaing vers, updating min/max as required
  for (let j = 1; j < verts.length; j++) {
    let temp = _vectorDotProduct(axis, verts[j]);
    if (temp < min) min = temp;
    if (temp > max) max = temp;
  }

  return { min: min, max: max };
}

function _getPerpendicularAxis(verts, index) {
  let pt1 = verts[index];
  let pt2 = index >= verts.length - 1 ? verts[0] : verts[index + 1]; // get the next index, or wrap around if at the end

  let axis = new SATPoint(-(pt2.y - pt1.y), pt2.x - pt1.x);
  axis.normalize();
  return axis;
}

function _vectorDotProduct(pt1, pt2) {
  return pt1.x * pt2.x + pt1.y * pt2.y;
}

function keyPressed() {
  if (key == "ArrowLeft") {
    shapeB.x -= 10;
  }

  if (key == "ArrowRight") {
    shapeB.x += 10;
  }

  if (key == "ArrowUp") {
    shapeB.y -= 10;
  }

  if (key == "ArrowDown") {
    shapeB.y += 10;
  }

  if (key == "r") {
    shapeB.rotation += 10;
  }
  //let ret = check();
  //console.log("result", ret);
}

function mouseClicked() {}

function getTransformedVerts(s) {
  let verts = [];

  let cX = 0;
  let cY = 0;

  let halfW = s.w / 2;
  let halfH = s.h / 2;

  verts.push(new SATPoint(cX + halfW, cY + halfH));
  verts.push(new SATPoint(cX - halfW, cY + halfH));
  verts.push(new SATPoint(cX - halfW, cY - halfH));
  verts.push(new SATPoint(cX + halfW, cY - halfH));

  let rotation = s.rotation;
  for (let i = 0; i < verts.length; i++) {
    if (rotation != 0) {
      let vert = verts[i];
      //console.log("vert bef:", vert);
      let hyp = Math.sqrt(Math.pow(vert.x, 2) + Math.pow(vert.y, 2));
      //console.log("Hyp:", hyp);
      let angle = Math.atan2(vert.y, vert.x);
      angle += rotation * (Math.PI / 180);
      vert.x = Math.cos(angle) * hyp;
      vert.y = Math.sin(angle) * hyp;
      //console.log("vert aft:", vert);
      verts[i] = vert;
    }
  }
  //console.log("vertices", verts);
  return verts;
}

function rval(tmp, rotation) {
  //let vert = { x: tmp.x + tmp.w / 2, y: tmp.y - tmp.h / 2 };
  let vert = { x: tmp.x, y: tmp.y };
  if (rotation != 0) {
    //console.log("vert bef:", vert);
    let hyp = Math.sqrt(Math.pow(vert.x, 2) + Math.pow(vert.y, 2));
    //console.log("Hyp:", hyp);
    let angle = Math.atan2(vert.y, vert.x);
    angle += rotation * (Math.PI / 180);
    vert.x = Math.cos(angle) * hyp;
    vert.y = Math.sin(angle) * hyp;
  }
  return vert;
}

class SATPoint {
  x = 0;
  y = 0;

  constructor(x = 0, y = 0) {
    this.x = x;
    this.y = y;
  }

  normalize() {
    this.magnitude = 1;
  }

  set magnitude(value) {
    let len = Math.sqrt(Math.pow(this.x, 2) + Math.pow(this.y, 2));
    if (len == 0) return;
    let ratio = value / len;
    this.x *= ratio;
    this.y *= ratio;
  }
  get magnitude() {
    return Math.sqrt(Math.pow(this.x, 2) + Math.pow(this.y, 2));
  }

  clone() {
    let clone = new SATPoint();
    clone.x = this.x;
    clone.y = this.y;
    return clone;
  }
}

function getVerts(s) {
  ver = [];
  let cX = 0;
  let cY = 0;

  let halfX = s.w / 2;
  let halfY = s.h / 2;

  ver.push(new SATPoint(cX - halfX, cY - halfY));
  ver.push(new SATPoint(cX + halfX, cY - halfY));
  ver.push(new SATPoint(cX + halfX, cY + halfY));
  ver.push(new SATPoint(cX - halfX, cY + halfY));

  /*
  ver.push({ x: s.x, y: s.y });
  ver.push({ x: s.x + s.w, y: s.y });
  ver.push({ x: s.x + s.w, y: s.y + s.h });
  ver.push({ x: s.x, y: s.y + s.h });
  */
  return ver;
}

class SATCollisionInfo {
  shapeA = null; // the first shape
  shapeB = null; // the second shape
  distance = 0; // how much overlap there is
  vector = new SATPoint(); // the direction you need to move - unit vector
  shapeAContained = false; // is object A contained in object B
  shapeBContained = false; // is object B contained in object A
  separation = new SATPoint(); // how far to separate
}

function _checkRangesForContainment(
  rangeA,
  rangeB,
  collisionInfo,
  flipResultPositions
) {
  if (flipResultPositions) {
    if (rangeA.max < rangeB.max || rangeA.min > rangeB.min)
      collisionInfo.shapeAContained = false;
    if (rangeB.max < rangeA.max || rangeB.min > rangeA.min)
      collisionInfo.shapeBContained = false;
  } else {
    if (rangeA.max > rangeB.max || rangeA.min < rangeB.min)
      collisionInfo.shapeAContained = false;
    if (rangeB.max > rangeA.max || rangeB.min < rangeA.min)
      collisionInfo.shapeBContained = false;
  }
}
