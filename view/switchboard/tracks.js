
const cellWidth = 20;
const cellHeight = 20;


function trackDiagramInit() {
    new TrackDiagram(80, 30);
}

/*****************************************************
 * 
 * TrackDiagram
 * 
 * The user interfaces which allows the user to interact
 * with the tracks. A TrackDiagram consists of a grid of
 * TrackDiagrammCell objects, which in turn are occupied
 * by tracks.
 * 
 *****************************************************/

class TrackDiagram {
    constructor(columns, rows) {
        this.columns = columns;
        this.rows = rows;
        this.grid = [];

        this.svg = document.getElementById("trackdiagram");
        this.svg.setAttribute("width", this.columns * cellWidth)
        this.svg.setAttribute("height", this.rows * cellHeight)

        for (var i = 0; i < columns; i++) {
            var line = document.createElementNS(svgNS, "line");
            line.setAttribute("y1", 0);
            line.setAttribute("y2", this.rows * cellHeight);
            line.setAttribute("x1", i * cellWidth);
            line.setAttribute("x2", i * cellWidth);
            line.classList.add("tracks-bgline");
            this.svg.appendChild(line);
        }

        for (var i = 0; i < rows; i ++) {
            var line = document.createElementNS(svgNS, "line");
            line.setAttribute("x1", 0);
            line.setAttribute("x2", this.columns * cellWidth);
            line.setAttribute("y1", i * cellHeight);
            line.setAttribute("y2", i * cellHeight);
            line.classList.add("tracks-bgline");
            this.svg.appendChild(line);
        }

        for (var x = 0; x < this.columns; x++) {
            var col = [];
            for (var y = 0; y < this.rows; y++) {
                col.push(new TrackDiagramCell(this, x, y))
            }
            this.grid.push(col)
        }
    }

    // Reconstructs a TrackDiagram from its JSON serialization.
    static deserialize(ser) {
        var dgrm = new TrackDiagram(ser.columns, ser.rows);
        
        // Deserialize all tracks and attach them to their cell
        for (var i in ser.tracks) {
            var obj = ser.tracks[i];
            var t;
            var cell = dgrm.grid[obj.c][obj.r];
            if (obj.kind <= switchHorizontalDiagonalLower) {
                t = new TrackSwitch(cell, obj.kind);
            } else if (obj.kind <= switchHorizontalCrossDiagonalBack) {
                t = new TrackSwitchCross(cell, obj.kind);
            } else if (obj.kind <= trackVerticalStopBottom) {
                t = new TrackSimple(cell, obj.kind);
            } else if (obj.kind <= trackVerticalBlock) {
                t = new TrackBlock(cell, obj.kind);
            } else {
                throw "Oooops";
            }
            // Call addTrack here, because multiple tracks can occupy a cell.
            cell.addTrack(t);
        }
    }
}

/*******************************************************
 * A TrackDiagram consists of a grid of TrackDiagramCells.
 * The user interacts with these cells.
 * A cell can be occupied by a track (or even two tracks
 * running diagonally in parallel).
 *******************************************************/

class TrackDiagramCell {
    constructor(diagram, x, y) {
        this.diagram = diagram;
        this.x = x;
        this.y = y;
        this.track = null;

        this.svg = document.createElementNS(svgNS, "g");
        this.diagram.svg.appendChild(this.svg);

        // Create hover effect
        var rect = document.createElementNS(svgNS, "rect");
        rect.setAttribute("x", x * cellWidth);
        rect.setAttribute("y", y * cellHeight);
        rect.setAttribute("width", cellWidth);
        rect.setAttribute("height", cellHeight);
        rect.classList.add("tracks-bgcell");
        this.diagram.svg.appendChild(rect);
        this.svgBG = rect;
    }

    // @param track is a Track object or an array of Track objects.
    setTrack(track) {
        if (this.track) {
            this.dropTrack();
        }
        this.track = track;
        if (Array.isArray(this.track)) {
            for (var t of this.track) {
                t.createSVG(this.svg);
            }
        } else {
            this.track.createSVG(this.svg);
        }
    }

    // Like setTrack, but this function can be invoked multiple times
    // in case multiple tracks occupy a cell. Used during deserialization.
    //
    // @param track is a Track object.
    addTrack(track) {
        if (!this.track) {
            this.setTrack(track);
            return
        }
        if (Array.isArray(this.track)) {
            this.track.push(track);
        } else {
            this.track = [this.track, track];
        }
        track.createSVG(this.svg);
    }

    // The TrackDiagrammCell above or null.
    nbTop() {
        if (this.y == 0) {
            return null;
        }
        return this.diagram.grid[this.x][this.y-1];
    }

    // The TrackDiagrammCell below or null.
    nbBottom() {
        if (this.y == this.diagram.rows - 1) {
            return null;
        }
        return this.diagram.grid[this.x][this.y+1];
    }

    // The TrackDiagrammCell to the left or null.
    nbLeft() {
        if (this.x == 0) {
            return null;
        }
        return this.diagram.grid[this.x-1][this.y];
    }

    // The TrackDiagrammCell to the right or null.
    nbRight() {
        if (this.x == this.diagram.columns - 1) {
            return null;
        }
        return this.diagram.grid[this.x+1][this.y];
    }

    // Removes the track(s), which currently occupy the cell.
    dropTrack() {
        if (!this.track) {
            return;
        }
        this._dropTrack(this.track);
    }

    // Internal helper function.
    // A Track can occupy multiple neighbouring cells.
    // Therefore, the function drops the track on those cells as well. 
    _dropTrack(track) {
        if (track != this.track) {
            return;
        }
        this.svg.innerHTML = "";
        this.conTop = false;
        this.conLeft = false;
        this.conRight = false;
        this.conBottom = false;
        this.track = null;
        if (this.nbTop()) {
            this.nbTop()._dropTrack(track);
        }
        if (this.nbBottom()) {
            this.nbBottom()._dropTrack(track);
        }
        if (this.nbLeft()) {
            this.nbLeft()._dropTrack(track);
        }
        if (this.nbRight()) {
            this.nbRight()._dropTrack(track);
        }
    }
}

/*****************************************************
 * 
 * Tracks
 * 
 * These classes derive from TrackElement and
 * represent tracks and their connection
 * independently of the cell-grid structures used in the UI.
 * The Track objects know how to draw themselves using SVG.
 * 
 *****************************************************/

const switchVerticalDiagonalUpper = 1;
const switchVerticalDiagonalBackUpper = 2;
const switchVerticalDiagonalBackLower = 3;
const switchVerticalDiagonalLower = 4;
const switchHorizontalDiagonalUpper = 5;
const switchHorizontalDiagonalBackUpper = 6;
const switchHorizontalDiagonalBackLower = 7;
const switchHorizontalDiagonalLower = 8;
const switchVerticalCrossDiagonal = 9;
const switchVerticalCrossDiagonalBack = 10;
const switchHorizontalCrossDiagonal = 11;
const switchHorizontalCrossDiagonalBack = 12;
const trackHorizontal = 20;
const trackVertical = 21;
const trackDiagonalUpper = 22;
const trackDiagonalLower = 23;
const trackDiagonalBackUpper = 24;
const trackDiagonalBackLower = 25;
const trackHorizontalStopRight = 26;
const trackHorizontalStopLeft = 27;
const trackVerticalStopTop = 28;
const trackVerticalStopBottom = 29;
const trackHorizontalBlock = 30;
const trackVerticalBlock = 31;

// Abstract base class.
class TrackElement {
    constructor(cell, kind) {
        this.cell = cell;
        this.kind = kind;
    }

    svgX() {
        return this.cell.x * cellWidth;
    }

    svgY() {
        return this.cell.y * cellHeight;
    }
}
  
// A vertical, horizontal or diagonal track.
class TrackSimple extends TrackElement {
    constructor(cell, kind) {
        super(cell, kind);
    }

    createSVG(svg) {
        this.svgParent = svg;

        if (this.kind == trackHorizontal) {
            var rect = document.createElementNS(svgNS, "rect");
            rect.setAttribute("x", this.svgX());
            rect.setAttribute("y", this.svgY() + cellHeight/3);
            rect.setAttribute("width", cellWidth);
            rect.setAttribute("height", cellHeight/3);
            rect.classList.add("track");
            this.svgParent.appendChild(rect);
        } else if (this.kind == trackVertical) {
            var rect = document.createElementNS(svgNS, "rect");
            rect.setAttribute("x", this.svgX() + cellWidth/3);
            rect.setAttribute("y", this.svgY());
            rect.setAttribute("width", cellWidth/3);
            rect.setAttribute("height", cellHeight);
            rect.classList.add("track");
            this.svgParent.appendChild(rect);
        } else if (this.kind == trackDiagonalUpper) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + this.svgX() + " " + (this.svgY() + cellHeight/3);
            d += " L" + (this.svgX() + cellWidth/3) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY());
            d += " L" + (this.svgX()) + " " + (this.svgY() + cellHeight*2/3)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == trackDiagonalLower) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX() + cellWidth/3) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight/3);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight*2/3);
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY() + cellHeight)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == trackDiagonalBackUpper) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX() + cellWidth/3) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight/3);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight*2/3)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == trackDiagonalBackLower) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX()) + " " + (this.svgY() + cellHeight/3);
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX() + cellWidth/3) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX()) + " " + (this.svgY() + cellHeight*2/3)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == trackHorizontalStopLeft) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX() + 2) + " " + (this.svgY() + cellHeight/6);
            d += " L" + (this.svgX()) + " " + (this.svgY() + cellHeight/6);
            d += " L" + (this.svgX()) + " " + (this.svgY() + cellHeight - cellHeight/6);
            d += " L" + (this.svgX() + 2) + " " + (this.svgY() + cellHeight - cellHeight/6);
            d += " L" + (this.svgX() + 2) + " " + (this.svgY() + cellHeight*2/3);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight*2/3);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight/3);
            d += " L" + (this.svgX() + 2) + " " + (this.svgY() + cellHeight/3)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == trackHorizontalStopRight) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX() + cellWidth - 2) + " " + (this.svgY() + cellHeight/6);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight/6);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight - cellHeight/6);
            d += " L" + (this.svgX() + cellWidth - 2) + " " + (this.svgY() + cellHeight - cellHeight/6);
            d += " L" + (this.svgX() + cellWidth - 2) + " " + (this.svgY() + cellHeight*2/3);
            d += " L" + (this.svgX()) + " " + (this.svgY() + cellHeight*2/3);
            d += " L" + (this.svgX()) + " " + (this.svgY() + cellHeight/3);
            d += " L" + (this.svgX() + cellWidth -2) + " " + (this.svgY() + cellHeight/3)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == trackVerticalStopTop) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX() + cellWidth/6) + " " + (this.svgY() + 2);
            d += " L" + (this.svgX() + cellWidth/6) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth - cellWidth/6) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth - cellWidth/6) + " " + (this.svgY() + 2);
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY() + 2);
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX() + cellWidth/3) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX() + cellWidth/3) + " " + (this.svgY() + 2)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == trackVerticalStopBottom) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX() + cellWidth/6) + " " + (this.svgY() + cellHeight - 2);
            d += " L" + (this.svgX() + cellWidth/6) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX() + cellWidth - cellWidth/6) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX() + cellWidth - cellWidth/6) + " " + (this.svgY() + cellHeight - 2);
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY() + cellHeight - 2);
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth/3) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth/3) + " " + (this.svgY() + cellHeight - 2)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else {
            throw "Oooops";
        }
    }
}
 
// A track with a single switch
class TrackSwitch extends TrackElement {
     constructor(cell, kind) {
         super(cell, kind);
         if (kind == switchHorizontalCrossDiagonal || kind == switchHorizontalCrossDiagonalBack || kind == switchVerticalCrossDiagonal || kind == switchVerticalCrossDiagonalBack) {
             throw "Ooops";
         }
     }

    createSVG(svg) {
        this.svgParent = svg;

        if (this.kind >= switchVerticalDiagonalUpper && this.kind <= switchVerticalDiagonalLower) {
            var rect = document.createElementNS(svgNS, "rect");
            rect.setAttribute("x", this.svgX() + cellWidth/3);
            rect.setAttribute("y", this.svgY());
            rect.setAttribute("width", cellWidth/3);
            rect.setAttribute("height", cellHeight);
            rect.classList.add("track");
            this.svgParent.appendChild(rect);
        } else {
            var rect = document.createElementNS(svgNS, "rect");
            rect.setAttribute("x", this.svgX());
            rect.setAttribute("y", this.svgY() + cellHeight/3);
            rect.setAttribute("width", cellWidth);
            rect.setAttribute("height", cellHeight/3);
            rect.classList.add("track");
            this.svgParent.appendChild(rect);
        }

        if (this.kind == switchVerticalDiagonalUpper || this.kind == switchHorizontalDiagonalUpper) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + this.svgX() + " " + (this.svgY() + cellHeight/3);
            d += " L" + (this.svgX() + cellWidth/3) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY());
            d += " L" + (this.svgX()) + " " + (this.svgY() + cellHeight*2/3)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == switchVerticalDiagonalLower || this.kind == switchHorizontalDiagonalLower) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX() + cellWidth/3) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight/3);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight*2/3);
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY() + cellHeight)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == switchVerticalDiagonalBackUpper || this.kind == switchHorizontalDiagonalBackUpper) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX() + cellWidth/3) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight/3);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight*2/3)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == switchVerticalDiagonalBackLower || this.kind == switchHorizontalDiagonalBackLower) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX()) + " " + (this.svgY() + cellHeight/3);
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX() + cellWidth/3) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX()) + " " + (this.svgY() + cellHeight*2/3)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else {
            throw "Oooops";
        }
    }

}

// A crossing track that can optionally switch.
class TrackSwitchCross extends TrackElement {
    constructor(cell, kind) {
        super(cell, kind);
        if (kind != switchHorizontalCrossDiagonal && kind != switchHorizontalCrossDiagonalBack && kind != switchVerticalCrossDiagonal && kind != switchVerticalCrossDiagonalBack) {
            throw "Ooops";
        }
    }

    createSVG(svg) {
        this.svgParent = svg;

        if (this.kind == switchHorizontalCrossDiagonal || this.kind == switchHorizontalCrossDiagonalBack) {
            var rect = document.createElementNS(svgNS, "rect");
            rect.setAttribute("x", this.svgX());
            rect.setAttribute("y", this.svgY() + cellHeight/3);
            rect.setAttribute("width", 2*cellWidth);
            rect.setAttribute("height", cellHeight/3);
            rect.classList.add("track");
            this.svgParent.appendChild(rect);
        } else {
            var rect = document.createElementNS(svgNS, "rect");
            rect.setAttribute("x", this.svgX() + cellWidth/3);
            rect.setAttribute("y", this.svgY());
            rect.setAttribute("width", cellWidth/3);
            rect.setAttribute("height", 2*cellHeight);
            rect.classList.add("track");
            this.svgParent.appendChild(rect);
        }
        if (this.kind == switchHorizontalCrossDiagonal) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX() + cellWidth/3) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX() + cellWidth + cellWidth*2/3) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth + cellWidth/3) + " " + (this.svgY())+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == switchHorizontalCrossDiagonalBack) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX() + cellWidth/3) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth*2/3) + " " + (this.svgY());
            d += " L" + (this.svgX() + cellWidth + cellWidth*2/3) + " " + (this.svgY() + cellHeight);
            d += " L" + (this.svgX() + cellWidth + cellWidth/3) + " " + (this.svgY() + cellHeight)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == switchVerticalCrossDiagonal) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX()) + " " + (this.svgY() + cellHeight + cellHeight/3);
            d += " L" + (this.svgX()) + " " + (this.svgY() + cellHeight + cellHeight*2/3);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight*2/3);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight/3)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else if (this.kind == switchVerticalCrossDiagonalBack) {
            var path = document.createElementNS(svgNS, "path");
            var d = "M" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight + cellHeight/3);
            d += " L" + (this.svgX() + cellWidth) + " " + (this.svgY() + cellHeight + cellHeight*2/3);
            d += " L" + (this.svgX()) + " " + (this.svgY() + cellHeight*2/3);
            d += " L" + (this.svgX()) + " " + (this.svgY() + cellHeight/3)+ " Z";
            path.setAttribute("d", d);
            path.classList.add("track");
            this.svgParent.appendChild(path);
        } else {
            throw "Oooops";
        }
    }
}
 
// A horizontal or vertical track that can display a train number.
class TrackBlock extends TrackElement {
    constructor(cell,kind) {
        super(cell, kind);
    }

    createSVG(svg) {
        this.svgParent = svg;

        if (this.kind == trackHorizontalBlock) {
            var rect = document.createElementNS(svgNS, "rect");
            rect.setAttribute("x", this.svgX());
            rect.setAttribute("y", this.svgY() + cellHeight/3);
            rect.setAttribute("width", cellWidth*4);
            rect.setAttribute("height", cellHeight/3);
            rect.classList.add("track");
            this.svgParent.appendChild(rect);

            rect = document.createElementNS(svgNS, "rect");
            rect.setAttribute("x", this.svgX() + cellWidth/2);
            rect.setAttribute("y", this.svgY() + cellHeight/6);
            rect.setAttribute("width", cellWidth*3);
            rect.setAttribute("height", cellHeight*2/3);
            rect.setAttribute("rx", 4);
            rect.setAttribute("ry", 4);
            rect.classList.add("track-block");
            this.svgParent.appendChild(rect);

            var text = document.createElementNS(svgNS, "text");
            text.setAttribute("x", this.svgX() + cellWidth*2);
            text.setAttribute("y", this.svgY() + cellHeight*2/3);
            text.setAttribute("text-anchor", "middle");
            this.cdata = document.createTextNode(this.train ? this.train.name : "");
            text.appendChild(this.cdata);
            text.classList.add("track-blocktext");
            this.svgParent.appendChild(text);
        } else if (this.kind == trackVerticalBlock) {
            var rect = document.createElementNS(svgNS, "rect");
            rect.setAttribute("x", this.svgX() + cellWidth/3);
            rect.setAttribute("y", this.svgY());
            rect.setAttribute("width", cellWidth/3);
            rect.setAttribute("height", cellHeight*4);
            rect.classList.add("track");
            this.svgParent.appendChild(rect);

            rect = document.createElementNS(svgNS, "rect");
            rect.setAttribute("x", this.svgX() + cellWidth/6);
            rect.setAttribute("y", this.svgY() + cellHeight/2);
            rect.setAttribute("width", cellWidth*2/3);
            rect.setAttribute("height", cellHeight*3);
            rect.setAttribute("rx", 4);
            rect.setAttribute("ry", 4);
            rect.classList.add("track-block");
            this.svgParent.appendChild(rect);

            var text = document.createElementNS(svgNS, "text");
            text.setAttribute("x", this.svgX() + cellWidth/3);
            text.setAttribute("y", this.svgY() + cellHeight*2);
            text.setAttribute("text-anchor", "middle");
            text.setAttribute("transform", "rotate(90, " + (this.svgX() + cellWidth/3) + ", " + (this.svgY() + cellHeight*2) + ")");
            this.cdata = document.createTextNode(this.train ? this.train.name : "");
            text.appendChild(this.cdata);
            text.classList.add("track-blocktext");
            this.svgParent.appendChild(text);
        }
    }

    setTrain(train) {
        if (this.train == train) {
            return;
        }
        this.train = train;
        if (!this.cdata) {
            return;
        }
        if (train) {
            this.cdata.nodeValue = train.name;
        } else {
            this.cdata.nodeValue = "";
        }
    }
 }
 
 /*****************************************************
 * 
 * Trains
 * 
 *****************************************************/

// Train is used by TrackBlock to indicate which train is occupying a block.
class Train {
    constructor(name) {
        this.name = name;
    }
}
