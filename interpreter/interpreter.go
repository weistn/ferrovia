package interpreter

import (
	"fmt"
	"strings"

	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model"
	"github.com/weistn/ferrovia/model/tracks"
	"github.com/weistn/ferrovia/parser"
)

type Interpreter struct {
	errlog           *errlog.ErrorLog
	ast              *parser.File
	model            *model.Model
	tracksWithAnchor []*tracks.Track
	identifiers      map[string]interface{}
}

func NewInterpreter(errlog *errlog.ErrorLog) *Interpreter {
	m := &model.Model{}
	m.Tracks = tracks.NewTrackSystem()
	return &Interpreter{errlog: errlog, model: m, identifiers: make(map[string]interface{})}
}

func (b *Interpreter) ProcessStatics(ast *parser.File) *model.Model {
	b.ast = ast

	// Determine all identifiers
	for _, s := range ast.Statements {
		if s == nil {
			break
		}
		switch t := s.(type) {
		case *parser.GroundPlate:
			// Do nothing by intention
		case *parser.Layer:
			if t.Name.StringValue != "" {
				b.identifiers[t.Name.StringValue] = s
			}
		case *parser.Switchboard:
			// Do nothing by intention
		case *parser.Tracks:
			if t.Name != nil {
				b.identifiers[t.Name.StringValue] = s
			}
		default:
			panic("Ooooops")
		}
	}

	// Compute all ground plates and layers
	for _, s := range ast.Statements {
		if s == nil {
			break
		}
		switch t := s.(type) {
		case *parser.GroundPlate:
			b.processGround(t)
		case *parser.Layer:
			b.processLayer(t)
		case *parser.Switchboard:
			b.processSwitchboard(t)
		case *parser.Tracks:
			// Do nothing by intention
		default:
			panic("Ooooops")
		}
	}

	// Compute all tracks
	for _, s := range ast.Statements {
		if s == nil {
			break
		}
		switch t := s.(type) {
		case *parser.GroundPlate:
			// Do nothing by intention
		case *parser.Layer:
			// Do nothing by intention
		case *parser.Switchboard:
			// Do nothing by intention
		case *parser.Tracks:
			if t.Name == nil {
				b.processTracks(t)
			}
		default:
			panic("Ooooops")
		}
	}

	// Determine the location of all tracks which are directly or indirectly anchored
	for _, t := range b.tracksWithAnchor {
		b.computeLocationFromAnchor(t)
	}

	// All tracks should be located by now
	for _, l := range b.model.Tracks.Layers {
		for _, t := range l.Tracks {
			if t.Location == nil {
				b.errlog.LogError(errlog.ErrorTracksWithoutPosition, t.SourceLocation)
				break
			}
		}
	}

	return b.model
}

func (b *Interpreter) processGround(ast *parser.GroundPlate) {
	ground := &model.GroundPlate{}
	ctx := NewGroundContext(ground)
	err := b.processStatements([]IContext{ctx}, ast.Expressions)
	if err != nil {
		return
	}
	err = ctx.Close(b)
	if err == nil {
		b.model.GroundPlates = append(b.model.GroundPlates, ground)
	}
}

func (b *Interpreter) processSwitchboard(ast *parser.Switchboard) {
	lines := strings.Split(ast.RawText, "\n")
	sb := processASCIIStructure(lines, ast.LocationText, b.errlog)
	if sb != nil {
		b.model.Switchboards = append(b.model.Switchboards, sb)
	}
}

func (b *Interpreter) processLayer(ast *parser.Layer) {
	if _, ok := b.model.Tracks.Layers[ast.Name.StringValue]; ok {
		b.errlog.LogError(errlog.ErrorDuplicateLayer, ast.Location, ast.Name.StringValue)
		return
	}
	ctx := NewLayerContext(&tracks.TrackLayer{Name: ast.Name.StringValue})
	err := b.processStatements([]IContext{ctx}, ast.Expressions)
	if err != nil {
		return
	}
	err = ctx.Close(b)
	if err == nil {
		b.model.Tracks.AddLayer(ctx.layer)
	}
}

func (b *Interpreter) processTracks(ast *parser.Tracks) {
	ctx := NewTracksContext(b.model.Tracks.Layers[""])
	err := b.processStatements([]IContext{ctx}, ast.Expressions)
	if err != nil {
		return
	}
	ctx.Close(b)
}

/*
func (b *Interpreter) processSwitchExpression(rail *parser.RailWay, exp *parser.Expression, columns []*railwayColumn, index int) (columnsNew []*railwayColumn, first *tracks.TrackConnection, last *tracks.TrackConnection) {

	l, ok := b.trackSystem.Layers[rail.Layer]
	if !ok {
		b.errlog.LogError(errlog.ErrorUnknownLayer, exp.Location, rail.Layer)
		return columns, columns[index].first, columns[index].last
	}
	newTrack := l.NewTrack(trackType)
	if newTrack == nil {
		b.errlog.LogError(errlog.ErrorUnknownTrackType, exp.Location, exp.Switch.Type)
		return columns, columns[index].first, columns[index].last
	}
	newTrack.SourceLocation = exp.Location

}
*/

/*
// `start` may be nil. In this case the resulting tracks form a block of tracks that has two unconnected ends.
// Otherwiese the new tracks will be connected to `start`.
func (b *Interpreter) processExpression(rail *parser.RailWay, exp *parser.Expression, start *tracks.TrackConnection) (first *tracks.TrackConnection, last *tracks.TrackConnection) {
	var c *tracks.TrackConnection
	if exp.Track != nil {
		c1, c2, err := b.processTrack(rail, exp, start)
		if err == nil {
			c = c1
			last = c2
		}
	} else if exp.Repeat != nil {
		c, last = b.processRepeatExpression(rail, exp, start)
	} else {
		panic("Oooops")
	}
	// This is the first track?
	if first == nil {
		first = c
	}
	return
}
*/

/*
func (b *Interpreter) processTrack(rail *parser.RailWay, ast *parser.Expression, prevCon *tracks.TrackConnection) (first *tracks.TrackConnection, last *tracks.TrackConnection, err *errlog.Error) {
	if ast.Track == nil {
		panic("Not a track")
	}
	l, ok := b.trackSystem.Layers[rail.Layer]
	if !ok {
		err = b.errlog.LogError(errlog.ErrorUnknownLayer, ast.Location, rail.Layer)
		return nil, nil, err
	}

	if named, ok := b.namedRailways[ast.Track.Type]; ok {
		if named.processed {
			err = b.errlog.LogError(errlog.ErrorNamedRailwayUsedTwice, ast.Location, ast.Track.Type)
			return nil, nil, err
		}
		b.processRailWay(named.ast)
		if named.column.first != nil && prevCon != nil {
			prevCon.Connect(named.column.first)
		}
		return named.column.first, named.column.last, nil
	}

	newTrack := l.NewTrack(ast.Track.Type)
	if newTrack == nil {
		err = b.errlog.LogError(errlog.ErrorUnknownTrackType, ast.Location, ast.Track.Type)
		return nil, nil, err
	}
	newTrack.SourceLocation = ast.Location
	if newTrack.Geometry.IncomingConnectionCount != 1 || newTrack.Geometry.OutgoingConnectionCount != 1 {
		panic("Track has more than one incoming or outgoing connections")
	}
	// Not a switch. Simply connect both ends
	first = newTrack.FirstConnection()
	if prevCon != nil {
		prevCon.Connect(newTrack.FirstConnection())
	}
	last = newTrack.SecondConnection()
	return
}
*/

/*
func (b *Interpreter) processTrackMark(exp *parser.Expression, track *tracks.Track) {
	if exp.TrackMark == nil {
		panic("Not a track")
	}
	if !track.AddMark(exp.TrackMark.Position, exp.TrackMark.Name) {
		b.errlog.LogError(errlog.ErrorTrackMarkDefinedTwice, exp.Location)
	}
}
*/

/*
func (b *Interpreter) processAnchor(exp *parser.Expression, con *tracks.TrackConnection) {
	if exp.Anchor == nil {
		panic("Oooops")
	}
	ast := exp.Anchor
	l := tracks.NewTrackLocation(con, tracks.Vec3{ast.X, ast.Y, ast.Z}, ast.Angle)
	if !con.Track.SetLocation(l) {
		b.errlog.LogError(errlog.ErrorTrackPositionedTwice, exp.Location)
	}
	b.tracksWithAnchor = append(b.tracksWithAnchor, con.Track)
}
*/

func (b *Interpreter) computeLocationFromAnchor(track *tracks.Track) {
	if track.Location == nil {
		panic("Track has no anchor")
	}
	tracks.NewEpoch()
	track.Tag()
	// println(track.Geometry.Name, track.Id, track.Location.Center[0], track.Location.Center[1], track.Location.Rotation)
	b.computeLocationOfConnectedTracks(track)
}

func (b *Interpreter) computeLocationOfConnectedTracks(track *tracks.Track) {
	for i := 0; i < track.ConnectionCount(); i++ {
		c := track.Connection(i)
		c2 := c.Opposite
		if c2 == nil || c2.Track.IsTagged() {
			continue
		}
		if c2.Track.Location != nil {
			b.errlog.LogError(errlog.ErrorTrackPositionedTwice, c2.Track.SourceLocation)
		}
		cpos, cangle := track.Location.Connection(i, track.Geometry)
		// println("    con at ", cpos[0], cpos[1], cangle, i)
		l := tracks.NewTrackLocation(c2, cpos, cangle)
		// println(c2.Track.Geometry.Name, c2.Track.Id, l.Center[0], l.Center[1], l.Rotation)
		c2.Track.SetLocation(l)
		c2.Track.Tag()
		b.computeLocationOfConnectedTracks(c2.Track)
	}
}

// The error returned (if any) is already logged. It just indicates that something went wrong
func (b *Interpreter) processStatements(ctx []IContext, ast []parser.IExpression) *errlog.Error {
	for _, exp := range ast {
		result, err := b.evalExpression(ctx, exp)
		if err != nil {
			return err
		}
		if result != nil {
			_, err = b.expandFunc(ctx, result, errlog.LocationRange{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Interpreter) lookup(ctx []IContext, loc errlog.LocationRange, name string) (*ExprValue, *errlog.Error) {
	for i := len(ctx) - 1; i >= 0; i-- {
		result, err := ctx[i].Lookup(b, loc, name)
		if err != nil || result != nil {
			return result, err
		}
	}
	return nil, b.errlog.LogError(errlog.ErrorUnknownMethod, loc, name)
}

func (b *Interpreter) evalExpression(ctx []IContext, expr parser.IExpression) (*ExprValue, *errlog.Error) {
	switch t := expr.(type) {
	case *parser.DotExpression:
		panic("TODO dot")
	case *parser.ContextExpression:
		newctx, err := b.evalToContext(ctx, t.Object)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Context %T\n", newctx)
		ctx = append(ctx, newctx)
		err = b.processStatements(ctx, t.Statements)
		if err != nil {
			return nil, err
		}
		return nil, newctx.Close(b)
	case *parser.CallExpression:
		f, err := b.evalExpression(ctx, t.Func)
		if err != nil {
			return nil, err
		}
		if f.Type != funcType {
			return nil, b.errlog.LogError(errlog.ErrorNotAMethod, errlog.LocationRange{})
		}
		return f.FuncValue.Func(b, ctx, errlog.LocationRange{}, t.Arguments...)
	case *parser.IdentifierExpression:
		ident, err := b.lookup(ctx, errlog.LocationRange{}, t.Identifier.StringValue)
		if err != nil {
			return nil, err
		}
		return ident, nil
	case *parser.BinaryExpression:
		return b.evalBinaryExpression(ctx, t)
	case *parser.DimensionExpression:
		return b.evalDimensionExpression(ctx, t)
	case *parser.ConstantExpression:
		result := &ExprValue{}
		if t.Value.Kind == parser.TokenInteger {
			result.Type = numberType
			// TODO: Check range
			result.NumberValue = float64(t.Value.IntegerValue.Int64())
		} else if t.Value.Kind == parser.TokenFloat {
			result.Type = numberType
			// TODO: Check range
			result.NumberValue, _ = t.Value.FloatValue.Float64()
		}
		return result, nil
	case *parser.VectorExpression:
		return b.evalVectorExpression(ctx, t)
	}
	fmt.Printf("Type %T\n", expr)
	panic("TODO")
}

func (b *Interpreter) evalVectorExpression(ctx []IContext, ast *parser.VectorExpression) (*ExprValue, *errlog.Error) {
	result := &ExprValue{Type: vectorType}
	for _, expr := range ast.Values {
		v, err := b.evalExpression(ctx, expr)
		if err != nil {
			return nil, err
		}
		result.VectorValue = append(result.VectorValue, v)
	}
	return result, nil
}

func (b *Interpreter) evalDimensionExpression(ctx []IContext, ast *parser.DimensionExpression) (*ExprValue, *errlog.Error) {
	left, err := b.evalExpression(ctx, ast.Value)
	if err != nil {
		return nil, err
	}
	switch ast.Dimension.StringValue {
	case "mm", "deg":
		return left, nil
	case "cm":
		val, err := b.ToFloat(left, errlog.LocationRange{})
		if err != nil {
			return nil, err
		}
		result := &ExprValue{Type: numberType}
		result.NumberValue = val * 10
		return result, nil
	case "m":
		val, err := b.ToFloat(left, errlog.LocationRange{})
		if err != nil {
			return nil, err
		}
		result := &ExprValue{Type: numberType}
		result.NumberValue = val * 1000
		return result, nil
	}
	panic("Oooops")
}

func (b *Interpreter) evalBinaryExpression(ctx []IContext, ast *parser.BinaryExpression) (*ExprValue, *errlog.Error) {
	left, err := b.evalExpression(ctx, ast.Left)
	if err != nil {
		return nil, err
	}
	switch ast.Op.Kind {
	case parser.TokenLogicalAnd:
		b, err := b.ToBool(left, ast.Op.Location)
		if err != nil {
			return nil, err
		}
		if !b {
			return &ExprValue{Type: numberType, NumberValue: 0}, nil
		}
	case parser.TokenLogicalOr:
		b, err := b.ToBool(left, ast.Op.Location)
		if err != nil {
			return nil, err
		}
		if b {
			return &ExprValue{Type: numberType, NumberValue: 1}, nil
		}
	}

	right, err := b.evalExpression(ctx, ast.Right)
	if err != nil {
		return nil, err
	}

	switch ast.Op.Kind {
	case parser.TokenLogicalAnd:
		return left.LogicalAnd(right, ast.Op.Location)
	case parser.TokenLogicalOr:
		return left.LogicalOr(right, ast.Op.Location)
	case parser.TokenEqual:
		return left.Equal(right, ast.Op.Location)
	case parser.TokenNotEqual:
		return left.NotEqual(right, ast.Op.Location)
	case parser.TokenLessOrEqual:
		return left.LessOrEqual(right, ast.Op.Location)
	case parser.TokenGreaterOrEqual:
		return left.GreaterOrEqual(right, ast.Op.Location)
	case parser.TokenLess:
		return left.Less(right, ast.Op.Location)
	case parser.TokenGreater:
		return left.Greater(right, ast.Op.Location)
	case parser.TokenPlus:
		return left.Plus(right, ast.Op.Location)
	case parser.TokenDash:
		return left.Minus(right, ast.Op.Location)
	case parser.TokenAsterisk:
		return left.Mul(right, ast.Op.Location)
	case parser.TokenSlash:
		return left.Div(right, ast.Op.Location)
	case parser.TokenPercent:
		return left.Rem(right, ast.Op.Location)
	case parser.TokenAmpersand:
		return left.BinaryAnd(right, ast.Op.Location)
	case parser.TokenPipe:
		return left.BinaryOr(right, ast.Op.Location)
	case parser.TokenCaret:
		return left.BinaryXor(right, ast.Op.Location)
	case parser.TokenBitClear:
		return left.BinaryAndNot(right, ast.Op.Location)
	case parser.TokenShiftLeft:
		return left.Lsh(right, ast.Op.Location)
	case parser.TokenShiftRight:
		return left.Rsh(right, ast.Op.Location)
	}
	panic("Oooops")
}

/*
func (b *Interpreter) call(loc errlog.LocationRange, ctx IContext, f *ExprValue, args []parser.IExpression) (*ExprValue, *errlog.Error) {
	if f.Type != funcType {
		return nil, errlog.NewError(errlog.ErrorNotAMethod, errlog.LocationRange{})
	}
	return f.FuncValue.Func(b, ctx, errlog.LocationRange{}, args)
}
*/

func (b *Interpreter) expandFunc(ctx []IContext, f *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if f.Type == funcType {
		return f.FuncValue.Func(b, ctx, loc)
	}
	return f, nil
}

func (b *Interpreter) ToFloat(e *ExprValue, loc errlog.LocationRange) (float64, *errlog.Error) {
	if e.Type == funcType {
		var err *errlog.Error
		e, err = e.FuncValue.Func(b, nil, errlog.LocationRange{})
		if err != nil {
			return 0, err
		}
	}
	if e.Type == numberType {
		return e.NumberValue, nil
	}
	return 0, b.errlog.LogError(errlog.ErrorTypeMismtach, loc)
}

func (b *Interpreter) ToBool(e *ExprValue, loc errlog.LocationRange) (bool, *errlog.Error) {
	if e.Type == funcType {
		var err *errlog.Error
		e, err = e.FuncValue.Func(b, nil, errlog.LocationRange{})
		if err != nil {
			return false, err
		}
	}
	if e.Type == numberType {
		return e.NumberValue != 0, nil
	}
	if e.Type == stringType {
		return e.StringValue != "", nil
	}
	if e.Type == vectorType {
		return len(e.VectorValue) != 0, nil
	}
	return false, b.errlog.LogError(errlog.ErrorTypeMismtach, loc)
}

func (b *Interpreter) ToVector(e *ExprValue, loc errlog.LocationRange) ([]*ExprValue, *errlog.Error) {
	if e.Type == funcType {
		var err *errlog.Error
		e, err = e.FuncValue.Func(b, nil, errlog.LocationRange{})
		if err != nil {
			return nil, err
		}
	}
	if e.Type == vectorType {
		return e.VectorValue, nil
	}
	return nil, b.errlog.LogError(errlog.ErrorTypeMismtach, loc)
}

func (b *Interpreter) ToString(e *ExprValue, loc errlog.LocationRange) (string, *errlog.Error) {
	if e.Type == funcType {
		var err *errlog.Error
		e, err = e.FuncValue.Func(b, nil, errlog.LocationRange{})
		if err != nil {
			return "", err
		}
	}
	if e.Type == stringType {
		return e.StringValue, nil
	}
	return "", b.errlog.LogError(errlog.ErrorTypeMismtach, loc)
}

func (b *Interpreter) ToContext(e *ExprValue, loc errlog.LocationRange) (IContext, *errlog.Error) {
	if e.Type == funcType {
		var err *errlog.Error
		e, err = e.FuncValue.Func(b, nil, errlog.LocationRange{})
		if err != nil {
			return nil, err
		}
	}
	if e.Type == contextType {
		return e.Context, nil
	}
	return nil, b.errlog.LogError(errlog.ErrorTypeMismtach, loc)
}

func (b *Interpreter) evalToFloat(ctx []IContext, expr parser.IExpression) (float64, *errlog.Error) {
	val, err := b.evalExpression(ctx, expr)
	if err != nil {
		return 0, err
	}
	return b.ToFloat(val, errlog.LocationRange{})
}

func (b *Interpreter) evalToString(ctx []IContext, expr parser.IExpression) (string, *errlog.Error) {
	val, err := b.evalExpression(ctx, expr)
	if err != nil {
		return "", err
	}
	return b.ToString(val, errlog.LocationRange{})
}

func (b *Interpreter) evalToContext(ctx []IContext, expr parser.IExpression) (IContext, *errlog.Error) {
	val, err := b.evalExpression(ctx, expr)
	if err != nil {
		return nil, err
	}
	return b.ToContext(val, errlog.LocationRange{})
}
