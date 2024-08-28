<script lang="ts">
  import { browser } from '$app/environment';
  import type { zPodSchema } from '$lib/models';
  import { pb } from '$lib/pocketbase';
  import type { Terminal } from '@xterm/xterm';
  import { onDestroy } from 'svelte';
  import type { z } from 'zod';

  export let pod: z.infer<typeof zPodSchema>;
  export let wantsConnectedToPod: boolean;

  let ws: WebSocket | undefined;
  let term: Terminal | undefined;
  let terminalDiv: HTMLElement;

  const initTerminal = async () => {
    const { Terminal } = await import('@xterm/xterm');

    term = new Terminal();
    term.open(terminalDiv);
  };

  const loadAttachAddon = async () => {
    const { AttachAddon } = await import('@xterm/addon-attach');

    if (term && ws) {
      term.loadAddon(new AttachAddon(ws));
      term.writeln(`connection to \x1B[1m${pod.name}\x1B[0m opened`);
    }
  };

  // --------------------------------------------------------------------------

  const connectWs = () => {
    const loc = window.location;
    ws = new WebSocket(
      `${loc.protocol === 'https:' ? 'wss' : 'ws'}://${loc.host}/api/noroom/pod/${pod.id}/attach?token=${pb.authStore.token}`,
    );

    ws.onopen = async () => {
      console.log('websocket open');

      loadAttachAddon();
    };

    ws.onclose = () => {
      console.log('websocket close');

      if (term) {
        term.writeln(`connection to \x1B[1m${pod.name}\x1B[0m lost`);
      }

      ws = undefined;
    };

    ws.onerror = (e) => {
      console.error('websocket error', e);
    };

    console.log('connectWs: done');
  };

  const disconnectWs = () => {
    if (ws) ws.close();
    ws = undefined;
  };

  // --------------------------------------------------------------------------

  const updateTerminalOpenState = (running: boolean) => {
    if (running && !ws) {
      connectWs();
      if (!term) {
        initTerminal();
      }
    }
  };

  $: if (browser) {
    if (pod.running) {
      if (wantsConnectedToPod) {
        updateTerminalOpenState(pod.running);
      } else {
        disconnectWs();
      }
    } else {
      wantsConnectedToPod = false;
    }
  }

  // --------------------------------------------------------------------------

  onDestroy(() => {
    disconnectWs();
  });
</script>

{#if pod.running}
  <div class="overflow-x-auto">
    <div bind:this={terminalDiv}></div>
  </div>
{/if}
