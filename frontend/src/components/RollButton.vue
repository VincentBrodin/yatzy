<script setup>
	import {ref} from 'vue';
	import socket from '@/lib/socket';

	const emit = defineEmits(["roll"]);

	const rolling = ref(false)

	socket.Emitter.addEventListener("message", (event) => {
		const callId = event.detail.callId;
		if (callId !== 1) return;

		console.log(event.detail.message.dice)
		emit("roll", event.detail.message.dice)

		setTimeout(() => {
			rolling.value = false;
		}, 750)
	});

	function roll() {
		rolling.value = true;
		socket.Send(1);
	}

</script>

<template>
	<button class="btn" :disabled="rolling" @click="roll">
		<template v-if="rolling">
			Rolling
		</template>
		<template v-else>
			Roll
		</template>
	</button>
</template>
