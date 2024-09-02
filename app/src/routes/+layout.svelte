<script lang="ts">
  import Navbar from '$lib/components/Navbar.svelte';
  import { pb } from '$lib/pocketbase';
  import { currentUser, shouldTrack } from '$lib/stores/user';
  import { getLastTrackedTime, sendTracking, setLastTrackedTime } from '$lib/tracking';
  import { onMount } from 'svelte';
  import '../app.css';

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

        $shouldTrack = false;
      }
    }
  };

  onMount(async () => {
    sendTrackingTimed();
  });

  $: if ($shouldTrack) sendTrackingTimed();
</script>

<div class="flex h-full w-full flex-col">
  <Navbar />

  <div class="grow overflow-y-auto">
    <slot />
  </div>
</div>
