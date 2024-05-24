<section>
	<h2>Parameters</h2>
	<div id="div_4306146210937178430">
		Region: <input type="number" class="TextBox" id="region" min="-1" max="255" /> <br />

		<table>
			<tbody>
				<tr title="Exchange stiffness (J/m)">
					<td>Aex</td>
					<td> <input type="text" class="TextBox" id="Aex" style="color: black;" /> J/m</td>
				</tr>

				<tr title="Landau-Lifshitz damping constant">
					<td>alpha</td>
					<td>
						<input type="text" class="TextBox" id="alpha" style="color: black;" />
					</td>
				</tr>

				<tr title="Cubic anisotropy direction #1">
					<td>anisC1</td>
					<td>
						<input type="text" class="TextBox" id="anisC1" style="color: black;" />
					</td>
				</tr>

				<tr title="Cubic anisotorpy directon #2">
					<td>anisC2</td>
					<td>
						<input type="text" class="TextBox" id="anisC2" style="color: black;" />
					</td>
				</tr>

				<tr title="Uniaxial anisotropy direction">
					<td>anisU</td>
					<td>
						<input type="text" class="TextBox" id="anisU" style="color: black;" />
					</td>
				</tr>

				<tr title="First magneto-elastic coupling constant (J/m3)">
					<td>B1</td>
					<td> <input type="text" class="TextBox" id="B1" style="color: black;" /> J/m3</td>
				</tr>

				<tr title="Second magneto-elastic coupling constant (J/m3)">
					<td>B2</td>
					<td> <input type="text" class="TextBox" id="B2" style="color: black;" /> J/m3</td>
				</tr>

				<tr title="Externally applied field (T)">
					<td>B_ext</td>
					<td> <input type="text" class="TextBox" id="B_ext" style="color: black;" /> T</td>
				</tr>

				<tr title="Bulk Dzyaloshinskii-Moriya strength (J/m2)">
					<td>Dbulk</td>
					<td> <input type="text" class="TextBox" id="Dbulk" style="color: black;" /> J/m2</td>
				</tr>

				<tr title="Interfacial Dzyaloshinskii-Moriya strength (J/m2)">
					<td>Dind</td>
					<td> <input type="text" class="TextBox" id="Dind" style="color: black;" /> J/m2</td>
				</tr>

				<tr title="Slonczewski secondairy STT term ε'">
					<td>EpsilonPrime</td>
					<td>
						<input type="text" class="TextBox" id="EpsilonPrime" style="color: black;" />
					</td>
				</tr>

				<tr title="exx component of the strain tensor">
					<td>exx</td>
					<td>
						<input type="text" class="TextBox" id="exx" style="color: black;" />
					</td>
				</tr>

				<tr title="exy component of the strain tensor">
					<td>exy</td>
					<td>
						<input type="text" class="TextBox" id="exy" style="color: black;" />
					</td>
				</tr>

				<tr title="exz component of the strain tensor">
					<td>exz</td>
					<td>
						<input type="text" class="TextBox" id="exz" style="color: black;" />
					</td>
				</tr>

				<tr title="eyy component of the strain tensor">
					<td>eyy</td>
					<td>
						<input type="text" class="TextBox" id="eyy" style="color: black;" />
					</td>
				</tr>

				<tr title="eyz component of the strain tensor">
					<td>eyz</td>
					<td>
						<input type="text" class="TextBox" id="eyz" style="color: black;" />
					</td>
				</tr>

				<tr title="ezz component of the strain tensor">
					<td>ezz</td>
					<td>
						<input type="text" class="TextBox" id="ezz" style="color: black;" />
					</td>
				</tr>

				<tr title="Slonczewski fixed layer polarization">
					<td>FixedLayer</td>
					<td>
						<input type="text" class="TextBox" id="FixedLayer" style="color: black;" />
					</td>
				</tr>

				<tr
					title="Slonczewski free layer thickness (if set to zero (default), then the thickness will be deduced from the mesh size) (m)"
				>
					<td>FreeLayerThickness</td>
					<td>
						<input type="text" class="TextBox" id="FreeLayerThickness" style="color: black;" /> m</td
					>
				</tr>

				<tr title="Defines spins that should be fixed">
					<td>frozenspins</td>
					<td>
						<input type="text" class="TextBox" id="frozenspins" style="color: black;" />
					</td>
				</tr>

				<tr title="Electrical current density (A/m2)">
					<td>J</td>
					<td> <input type="text" class="TextBox" id="J" style="color: black;" /> A/m2</td>
				</tr>

				<tr title="1st order cubic anisotropy constant (J/m3)">
					<td>Kc1</td>
					<td> <input type="text" class="TextBox" id="Kc1" style="color: black;" /> J/m3</td>
				</tr>

				<tr title="2nd order cubic anisotropy constant (J/m3)">
					<td>Kc2</td>
					<td> <input type="text" class="TextBox" id="Kc2" style="color: black;" /> J/m3</td>
				</tr>

				<tr title="3rd order cubic anisotropy constant (J/m3)">
					<td>Kc3</td>
					<td> <input type="text" class="TextBox" id="Kc3" style="color: black;" /> J/m3</td>
				</tr>

				<tr title="1st order uniaxial anisotropy constant (J/m3)">
					<td>Ku1</td>
					<td> <input type="text" class="TextBox" id="Ku1" style="color: black;" /> J/m3</td>
				</tr>

				<tr title="2nd order uniaxial anisotropy constant (J/m3)">
					<td>Ku2</td>
					<td> <input type="text" class="TextBox" id="Ku2" style="color: black;" /> J/m3</td>
				</tr>

				<tr title="Slonczewski Λ parameter">
					<td>Lambda</td>
					<td>
						<input type="text" class="TextBox" id="Lambda" style="color: black;" />
					</td>
				</tr>

				<tr title="Height of vertically magnetized part of MFM tip">
					<td>MFMDipole</td>
					<td> <input type="text" class="TextBox" id="MFMDipole" style="color: black;" /> m</td>
				</tr>

				<tr title="MFM lift height">
					<td>MFMLift</td>
					<td> <input type="text" class="TextBox" id="MFMLift" style="color: black;" /> m</td>
				</tr>

				<tr title="Saturation magnetization (A/m)">
					<td>Msat</td>
					<td> <input type="text" class="TextBox" id="Msat" style="color: black;" /> A/m</td>
				</tr>

				<tr
					title="Disable magnetostatic interaction per region (default=0, set to 1 to disable). E.g.: NoDemagSpins.SetRegion(5, 1) disables the magnetostatic interaction in region 5."
				>
					<td>NoDemagSpins</td>
					<td>
						<input type="text" class="TextBox" id="NoDemagSpins" style="color: black;" />
					</td>
				</tr>

				<tr title="Electrical current polarization">
					<td>Pol</td>
					<td>
						<input type="text" class="TextBox" id="Pol" style="color: black;" />
					</td>
				</tr>

				<tr title="Temperature (K)">
					<td>Temp</td>
					<td> <input type="text" class="TextBox" id="Temp" style="color: black;" /> K</td>
				</tr>

				<tr title="Non-adiabaticity of spin-transfer-torque">
					<td>xi</td>
					<td>
						<input type="text" class="TextBox" id="xi" style="color: black;" />
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</section>

<style>
	section {
		grid-area: parameters;
	}
</style>
