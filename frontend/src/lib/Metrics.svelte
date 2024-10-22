<script lang="ts">
	import { metricsState as m } from '$api/incoming/metrics';
	import { postResetError } from '$api/outgoing/metrics';
	import { Card, Progressbar, Button } from 'flowbite-svelte';
</script>

<section>
	<h2 class="mb-4 text-2xl font-semibold">Metrics</h2>
	{#if $m.error}
		<p class="text-red-500">Got an error while collecting system metrics: {$m.error}</p>
		<Button on:click={postResetError} color="red" class="mt-4" outline>Retry</Button>
	{:else}
		<div class="grid gap-6 md:grid-cols-2">
			<!-- System Metrics Card -->
			<Card class="p-6">
				<h3 class="mb-4 text-xl font-semibold">System Metrics</h3>
				<div class="space-y-6">
					<!-- Global System Metrics -->
					<Card class="p-4">
						<h4 class="mb-2 text-lg font-semibold">Global</h4>
						<ul class="space-y-4">
							<li>
								<div class="flex items-center justify-between">
									<span class="font-medium">CPU Usage:</span>
									<span>{$m.cpuPercentTotal.toFixed(2)}%</span>
								</div>
								<Progressbar animate progress={$m.cpuPercentTotal} color="blue" />
							</li>
							<li>
								<div class="flex items-center justify-between">
									<span class="font-medium">RAM Usage:</span>
									<span>{$m.ramPercentTotal.toFixed(2)}%</span>
								</div>
								<Progressbar animate progress={$m.ramPercentTotal} color="green" />
							</li>
						</ul>
					</Card>
					<!-- Simulation System Metrics -->
					<Card class="p-4">
						<h4 class="mb-2 text-lg font-semibold">Simulation</h4>
						<ul class="space-y-4">
							<li>
								<span class="font-medium">PID:</span>
								{$m.pid}
							</li>
							<li>
								<div class="flex items-center justify-between">
									<span class="font-medium">CPU Usage:</span>
									<span>{$m.cpuPercent.toFixed(2)}%</span>
								</div>
								<Progressbar animate progress={$m.cpuPercent} color="blue" />
							</li>
							<li>
								<div class="flex items-center justify-between">
									<span class="font-medium">RAM Usage:</span>
									<span>{$m.ramPercent.toFixed(2)}%</span>
								</div>
								<Progressbar animate progress={$m.ramPercent} color="green" />
							</li>
						</ul>
					</Card>
				</div>
			</Card>

			<!-- GPU Metrics Card -->
			<Card class="p-6">
				<h3 class="mb-4 text-xl font-semibold">GPU Metrics</h3>
				<div class="space-y-6">
					<!-- Global GPU Metrics -->
					<Card class="p-4">
						<h4 class="font- mb-2 align-middle text-lg">Global</h4>
						<ul class="space-y-4">
							<li>
								<span class="font-medium">Name:</span>
								{$m.gpuName}
							</li>
							<li>
								<span class="font-medium">Temperature:</span>
								{$m.gpuTemperature}°C
							</li>
							<li>
								<div class="flex items-center justify-between">
									<span class="font-medium">Utilization:</span>
									<span>{$m.gpuUtilizationPercent}%</span>
								</div>
								<Progressbar animate progress={$m.gpuUtilizationPercent} color="purple" />
							</li>
							<li>
								<div class="flex items-center justify-between">
									<span class="font-medium">Power Draw:</span>
									<span>{$m.gpuPowerDraw.toFixed(2)} W</span>
								</div>
								<Progressbar
									animate
									progress={($m.gpuPowerDraw / $m.gpuPowerLimit) * 100}
									color="purple"
								/>
							</li>
							<li>
								<div class="flex items-center justify-between">
									<span class="font-medium">VRAM Total:</span>
									<span>{($m.gpuVramTotal / 1024).toFixed(2)} GiB</span>
								</div>
							</li>
						</ul>
					</Card>
					<!-- Simulation GPU Metrics -->
					<Card class="p-4">
						<h4 class="mb-2 text-lg font-semibold">Simulation</h4>
						<ul class="space-y-4">
							<li>
								<div class="flex items-center justify-between">
									<span class="font-medium">VRAM Used:</span>
									<span>{($m.gpuVramUsed / 1024).toFixed(2)} GiB</span>
								</div>
								<Progressbar
									animate
									progress={($m.gpuVramUsed / $m.gpuVramTotal) * 100}
									color="purple"
								/>
							</li>
						</ul>
					</Card>
				</div>
			</Card>
		</div>
	{/if}
</section>

<style>
	section {
		grid-area: metrics;
	}
</style>
