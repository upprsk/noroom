<script lang="ts">
  import Navbar from '$lib/components/Navbar.svelte';
  import type { UserModel } from '$lib/models';
  import { currentUser } from '$lib/stores/user';
  import { onMount } from 'svelte';
  import '../app.css';
  import { sendFingerprint } from '$lib/pocketbase';

  export let data;

  // Set the current user from the data passed in from the server
  $: currentUser.set(data.user as UserModel);

  onMount(async () => {
    const { getFingerprint, getFingerprintData } = await import('@thumbmarkjs/thumbmarkjs');
    const fingerprint = await getFingerprint();
    if (typeof fingerprint === 'string') {
      sendFingerprint($currentUser?.id, fingerprint);
    }

    const data = await getFingerprintData();
    console.log(data);
  });
</script>

<div class="flex flex-col w-full h-full">
  <Navbar />

  <div class="grow overflow-y-auto">
    <slot />
  </div>
</div>
