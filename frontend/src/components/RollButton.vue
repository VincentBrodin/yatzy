<script setup>
	import socket from '@/lib/socket';

	const emit = defineEmits(["roll"]);
	const props = defineProps({
		rolling: Boolean
	});

	function roll() {
		emit("roll")
		socket.Send(1);
	}

</script>

<template>
	<button class="btn-lg btn-primary w-full" :disabled="props.rolling" @click="roll">
		<template v-if="props.rolling">
			Rolling
		</template>
		<template v-else>
			Roll
		</template>
	</button>
</template>
