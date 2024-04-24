<div class="solver">
	<span
		title="Click to show/hide"
		style="cursor:pointer; font-size:1.2em; font-weight:bold; color:gray">▾ solver</span
	> <br />
	<div id="div_359618236257259938">
		Type: <select id="solvertype">
			<option value="bw_euler"> bw_euler </option>
			<option value="euler"> euler </option>
			<option value="heun"> heun </option>
			<option value="rk4"> rk4 </option>
			<option value="rk23"> rk23 </option>
			<option value="rk45"> rk45 </option>
			<option value="rkf56"> rkf56 </option>
		</select>
		<table>
			<tbody>
				<tr>
					<td>
						<table>
							<tbody>
								<tr title="">
									<td> <button id="run" class="Button">Run</button></td>
									<td>
										<input
											type="text"
											class="TextBox"
											id="runtime"
											size="8"
											style="color: black;"
										/>s</td
									>
								</tr>
								<tr title="">
									<td> <button id="steps" class="Button">Steps</button></td>
									<td>
										<input
											type="text"
											class="TextBox"
											id="runsteps"
											size="8"
											style="color: black;"
										/>
									</td>
								</tr>
								<tr title="">
									<td> <button id="relax" class="Button">Relax</button></td>
									<td></td>
								</tr>
								<tr title="Break current run loop">
									<td> <button id="break" class="Button">Break</button></td>
									<td></td>
								</tr>
							</tbody>
						</table>
					</td>
					<td> &nbsp; &nbsp; &nbsp; &nbsp; </td>
					<td>
						<table>
							<tbody>
								<tr title="Time steps taken">
									<td>step: </td>
									<td><span id="nsteps">13149</span> </td>
								</tr>
								<tr title="">
									<td>time: </td>
									<td><span id="time">1.00000e-09</span> s </td>
								</tr>
								<tr title="">
									<td>dt: </td>
									<td><span id="dt">4.394e-14</span> s </td>
								</tr>
								<tr title="Maximum relative error/step">
									<td>err/step: </td>
									<td><span id="lasterr">1.991e-09</span> </td>
								</tr>
								<tr title="Maximum absolute torque">
									<td>MaxTorque:</td>
									<td><span id="maxtorque">6.667e+00 T</span> </td>
								</tr>
							</tbody>
						</table>
					</td>
					<td> &nbsp; &nbsp; &nbsp; &nbsp; </td>
					<td>
						<table>
							<tbody>
								<tr title="">
									<td>fixdt: </td>
									<td
										><input type="text" class="TextBox" id="fixdt" size="8" style="color: black;" />
										s
									</td>
								</tr>
								<tr title="">
									<td>mindt: </td>
									<td
										><input type="text" class="TextBox" id="mindt" size="8" style="color: black;" />
										s
									</td>
								</tr>
								<tr title="">
									<td>maxdt: </td>
									<td
										><input type="text" class="TextBox" id="maxdt" size="8" style="color: black;" />
										s
									</td>
								</tr>
								<tr title="">
									<td>maxerr: </td>
									<td
										><input
											type="text"
											class="TextBox"
											id="maxerr"
											size="8"
											style="color: black;"
										/>/step</td
									>
								</tr>
							</tbody>
						</table>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>
