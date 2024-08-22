<script lang="ts">
  import Navbar from '$lib/components/Navbar.svelte';
  import type { UserModel } from '$lib/models';
  import { pb } from '$lib/pocketbase';
  import { currentUser } from '$lib/stores/user';
  import { getLastTrackedTime, sendTracking, setLastTrackedTime } from '$lib/tracking';
  import { onMount } from 'svelte';
  import '../app.css';

  export let data;

  // Set the current user from the data passed in from the server
  $: currentUser.set(data.user as UserModel);

  const checkTimeElapsed = () => {
    const ltt = getLastTrackedTime();
    if (!isNaN(ltt.valueOf())) {
      const delta = new Date().getTime() - ltt.getTime();

      const seconds = 1000;
      const minutes = seconds * 60;
      const hours = minutes * 60;

      // only update at most 1 time per hour
      if (delta / hours < 1) return false;
    }

    return true;
  };

  const sendTrackingTimed = async () => {
    if (checkTimeElapsed()) {
      const { getFingerprint, getFingerprintData } = await import('@thumbmarkjs/thumbmarkjs');
      const fingerprint = await getFingerprint();
      const data = await getFingerprintData();

      if (typeof fingerprint === 'string') {
        if ($currentUser) setLastTrackedTime();
        sendTracking(pb, $currentUser?.id, fingerprint, data);
      }
    }
  };

  onMount(async () => {
    sendTrackingTimed();
  });
</script>

<div class="flex flex-col w-full h-full">
  <Navbar />

  <div class="grow overflow-y-auto">
    <slot />
  </div>
</div>
