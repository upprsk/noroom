<script lang="ts">
  export let name: string;
  export let errors: string[] | undefined;
  export let value: string;
  export let constraints: { [key: string]: unknown } | undefined;

  export let disabled = false;
  export let placeholder = name;

  type TheEvent = KeyboardEvent & {
    currentTarget: EventTarget & HTMLTextAreaElement;
  };

  const onKeydown = (e: TheEvent) => {
    if (e.key == 'Tab') {
      e.preventDefault();

      const start = e.currentTarget.selectionStart;
      const end = e.currentTarget.selectionEnd;

      const tab = '    ';
      e.currentTarget.value =
        e.currentTarget.value.substring(0, start) + tab + e.currentTarget.value.substring(end);
      e.currentTarget.selectionStart = e.currentTarget.selectionEnd = start + tab.length;
    }
  };
</script>

<label class="form-control w-full max-w-xl">
  <div class="label">
    <span class="label-text"><slot /></span>
  </div>
  <textarea
    class="textarea textarea-bordered w-full max-w-xl min-h-52"
    on:keydown={onKeydown}
    aria-invalid={errors ? 'true' : undefined}
    bind:value
    {name}
    {disabled}
    {placeholder}
    {...constraints}
  />
  <div class="label">
    {#if errors}<span class="label-text-alt text-error">{errors}</span>{/if}
  </div>
</label>
