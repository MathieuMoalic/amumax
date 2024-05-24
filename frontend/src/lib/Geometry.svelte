<section>
	<h2>Geometry</h2>
	<div id="div_4350813584396622477">
		SetGeom( <select id="geomselect">
			<option value="Cell"> Cell </option>
			<option value="Circle"> Circle </option>
			<option value="Cone"> Cone </option>
			<option value="Cuboid"> Cuboid </option>
			<option value="Cylinder"> Cylinder </option>
			<option value="Diamond"> Diamond </option>
			<option value="Ellipse"> Ellipse </option>
			<option value="Ellipsoid"> Ellipsoid </option>
			<option value="GrainRoughness"> GrainRoughness </option>
			<option value="Hexagon"> Hexagon </option>
			<option value="ImageShape"> ImageShape </option>
			<option value="Layer"> Layer </option>
			<option value="Layers"> Layers </option>
			<option value="Rect"> Rect </option>
			<option value="RTriangle"> RTriangle </option>
			<option value="ShapeFromRegion"> ShapeFromRegion </option>
			<option value="Square"> Square </option>
			<option value="Squircle"> Squircle </option>
			<option value="Triangle"> Triangle </option>
			<option value="Universe"> Universe </option>
			<option value="Wave"> Wave </option>
			<option value="XRange"> XRange </option>
			<option value="YRange"> YRange </option>
			<option value="ZRange"> ZRange </option>
		</select>
		<input type="text" class="TextBox" id="geomargs" style="color: black;" />
		) <button id="setgeom" class="Button">Set</button>
		<br />
		<span id="geomdoc" style="color:gray"></span>
	</div>
</section>

<style>
	section {
		grid-area: geometry;
	}
</style>
