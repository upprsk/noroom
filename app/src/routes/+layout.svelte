<script lang="ts">
  import Navbar from '$lib/components/Navbar.svelte';
  import type { UserModel } from '$lib/models';
  import { currentUser } from '$lib/stores/user';
  import { sendTracking } from '$lib/tracking';
  import { onMount } from 'svelte';
  import '../app.css';

  export let data;

  // Set the current user from the data passed in from the server
  $: currentUser.set(data.user as UserModel);

  onMount(async () => {
    const { getFingerprint, getFingerprintData } = await import('@thumbmarkjs/thumbmarkjs');
    const fingerprint = await getFingerprint();
    const data = await getFingerprintData();

    if (typeof fingerprint === 'string') {
      sendTracking($currentUser?.id, fingerprint, data);
    }
  });
</script>

<div class="flex flex-col w-full h-full">
  <Navbar />

  <div class="grow overflow-y-auto">
    <slot />
  </div>
</div>
