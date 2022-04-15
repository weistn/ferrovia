
function renderCanvas(c, svgView, svgMeasure, svgGround, svgTracks) {
    document.getElementById("title").innerText = c.name;

    svgMeasure.innerHTML = "";
    svgGround.innerHTML = "";
    svgTracks.innerHTML = "";
    
    // TODO: hard-coded zoom factor
    svgView.setAttribute("width", c.width/2 + 10 + "px");
    svgView.setAttribute("height", c.height/2 + 10 + "px");

    // Render grid
    for (var i = 1; i * 100 < c.width; i++) {
        var svgText = document.createElementNS(svgNS, "text");
        svgText.classList.add("view2d-grid-measure");
        svgText.setAttribute("x", i * 100 + 8);
        svgText.setAttribute("y", 16);
        // TODO: hard-coded zoom factor
        svgText.innerHTML = (20*i).toString();
        svgMeasure.appendChild(svgText);
    }

    for (var i = 1; i * 100 < c.height; i++) {
        var svgText = document.createElementNS(svgNS, "text");
        svgText.classList.add("view2d-grid-measure");
        svgText.setAttribute("y", i * 100 + 16);
        svgText.setAttribute("x", 8);
        // TODO: hard-coded zoom factor
        svgText.innerHTML = (20*i).toString();
        svgMeasure.appendChild(svgText);
    }

    for (ground of c.ground) {
        if (ground.polygon) {
            var svgPoly = document.createElementNS(svgNS, "polygon");
            var points = [];
            for (var p of ground.polygon) {
                points.push((ground.left + p.x).toString() + "," + (ground.top + p.y).toString())
            }
            svgPoly.setAttribute("points", points.join(" "));
            svgPoly.classList.add("ground");
            svgGround.appendChild(svgPoly);
        } else {
            var svgRect = document.createElementNS(svgNS, "rect");
            svgRect.setAttribute("x", ground.left);
            svgRect.setAttribute("y", ground.top);
            svgRect.setAttribute("width", ground.width);
            svgRect.setAttribute("height", ground.height);
            svgRect.classList.add("ground");
            svgGround.appendChild(svgRect);
        }
    }

    if (!c.layers) {
        c.layers = [];
    }
    for (layer of c.layers) {
        for (track of layer.tracks) {
            if (track.l) {
                for (var line of track.l) {
                    var svgLine = document.createElementNS(svgNS, "path");
                    svgLine.classList.add("track-bars");
                    var sin = Math.sin(line.a * Math.PI/180);
                    var cos = Math.cos(line.a * Math.PI/180);    
                    var dx = sin * line.l;
                    var dy = -cos * line.l;        
                    var d = [
                        "M", line.x, line.y,
                        "L", line.x + dx, line.y + dy
                    ].join(" ");
                    svgLine.setAttribute("d", d);
                    svgTracks.appendChild(svgLine);

                    svgLine = document.createElementNS(svgNS, "path");
                    svgLine.classList.add("track-iron");
                    offsetY = -sin * 16.4 / 2;
                    offsetX = -cos * 16.4 / 2;        
                    d = [
                        "M", line.x + offsetX, line.y + offsetY,
                        "L", line.x + offsetX + dx, line.y + offsetY + dy,
                        "M", line.x - offsetX, line.y - offsetY,
                        "L", line.x - offsetX + dx, line.y - offsetY + dy
                    ].join(" ");
                    svgLine.setAttribute("d", d);
                    svgTracks.appendChild(svgLine);
                }
            }
            if (track.a) {
                for (var arc of track.a) {
                    var svgLine = document.createElementNS(svgNS, "path");
                    svgLine.classList.add("track-bars");
                    var d = describeArc(arc.cx, arc.cy, arc.r, arc.sa, arc.sa + arc.ta);
                    svgLine.setAttribute("d", d);
                    svgTracks.appendChild(svgLine);

                    svgLine = document.createElementNS(svgNS, "path");
                    svgLine.classList.add("track-iron");
                    var d1 = describeArc(arc.cx, arc.cy, arc.r - 16.4 / 2, arc.sa, arc.sa + arc.ta);
                    var d2 = describeArc(arc.cx, arc.cy, arc.r + 16.4 / 2, arc.sa, arc.sa + arc.ta);
                    svgLine.setAttribute("d", d1 + " " + d2);
                    svgTracks.appendChild(svgLine);            
                }
            }
            if (track.d) {
                var svgLine = document.createElementNS(svgNS, "path");
                svgLine.classList.add("track-delimiter");
                var d = "";
                for (var delimiter of track.d) {
                    d += [
                        " M", delimiter.x1, delimiter.y1,
                        "L", delimiter.x2, delimiter.y2
                    ].join(" ");
                }
                svgLine.setAttribute("d", d);
                svgTracks.appendChild(svgLine);
            }
        }
    }
}

function polarToCartesian(centerX, centerY, radius, angleInDegrees) {
    var angleInRadians = (angleInDegrees-90) * Math.PI / 180.0;
  
    return {
      x: centerX + (radius * Math.cos(angleInRadians)),
      y: centerY + (radius * Math.sin(angleInRadians))
    };
}
  
function describeArc(x, y, radius, startAngle, endAngle) {  
    var start = polarToCartesian(x, y, radius, endAngle);
    var end = polarToCartesian(x, y, radius, startAngle);
  
    var largeArcFlag = endAngle - startAngle <= 180 ? "0" : "1";
  
    var d = [
        "M", start.x, start.y, 
        "A", radius, radius, 0, largeArcFlag, 0, end.x, end.y
    ].join(" ");
  
    return d;       
}
