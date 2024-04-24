<span
	title="Click to show/hide"
	style="cursor:pointer; font-size:1.2em; font-weight:bold; color:gray">▾ mesh</span
> <br />
<div id="div_4506646511619690965">
	<table>
		<tbody>
			<tr>
				<td>gridsize: </td>
				<td><input type="text" class="TextBox" id="nx" size="8" style="color: black;" /> </td>
				<td> × <input type="text" class="TextBox" id="ny" size="8" style="color: black;" /></td>
				<td> × <input type="text" class="TextBox" id="nz" size="8" style="color: black;" /></td>
				<td> cells </td>
			</tr>
			<tr>
				<td>cellsize: </td>
				<td><input type="text" class="TextBox" id="cx" size="8" style="color: black;" /> </td>
				<td> × <input type="text" class="TextBox" id="cy" size="8" style="color: black;" /></td>
				<td> × <input type="text" class="TextBox" id="cz" size="8" style="color: black;" /></td>
				<td> m<sup>3</sup> </td>
			</tr>
			<tr>
				<td>PBC: </td>
				<td><input type="text" class="TextBox" id="px" size="8" style="color: black;" /> </td>
				<td> × <input type="text" class="TextBox" id="py" size="8" style="color: black;" /></td>
				<td> × <input type="text" class="TextBox" id="pz" size="8" style="color: black;" /></td>
				<td> repetitions </td>
			</tr>
			<tr>
				<td>worldsize:</td>
				<td><span id="wx">300</span> </td>
				<td> × <span id="wy">300</span></td>
				<td> × <span id="wz">20</span></td>
				<td> nm<sup>3</sup> </td>
			</tr>
		</tbody>
	</table>

	<button id="setmesh" class="Button">update</button>
	<span id="setmeshwarn"></span>
</div>
