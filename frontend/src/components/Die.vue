<script setup>
	import {ref} from 'vue';

	const props = defineProps({
		target: Number
	});

	const preview = ref(1);
	const rolling = ref(false);
	const selected = ref(false);

	function getRandomInt(max) {
		return Math.floor(Math.random() * max);
	}


	let rollInterval = null
	function roll() {
		if (selected.value) return
		rolling.value = true
		if (rollInterval) clearInterval(rollInterval)

		rollInterval = setInterval(() => {
			preview.value = getRandomInt(6) + 1
		}, 100)

		setTimeout(() => {
			clearInterval(rollInterval)
			rollInterval = null
			preview.value = props.target
			rolling.value = false
		}, 750)
	}

	function select() {
		selected.value = !selected.value
	}
	defineExpose({roll})
</script>

<template>
	<div class="die-container">
		<button @click="select">
			<img class="w-full h-auto transition-opacity" :class="{ 'opacity-50': rolling || selected  }"
				:src="`/img/${preview}.png`" alt="Dice Face">
		</button>
	</div>
</template>

<style scoped>
	.die-container {
		max-width: 100%;
		aspect-ratio: 1 / 1;
	}
</style>
