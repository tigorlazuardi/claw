<script lang="ts">
  import { getSourceServiceClient } from "../../connectrpc";
  import { toDate } from "../../connectrpc/js_date";
  import { parseCronExpression } from "cron-schedule";
  import { Popover, Combobox } from "bits-ui";
  import IconInfo from "@lucide/svelte/icons/info";
  import { resource, watch } from "runed";
  import { theme } from "../../store/theme";

  interface Props {
    value: string[];
  }

  let scheduleInputValue = $state("");

  watch(
    () => scheduleInputValue,
    (val) => {
      console.log($state.snapshot(val));
    },
  );

  const scheduleInputNextRunResource = resource(
    () => scheduleInputValue,
    async (pattern, _, { signal }) => {
      if (!pattern.trim()) {
        return;
      }
      return getNextRun(pattern, signal);
    },
    {
      debounce: 300,
    },
  );

  async function getNextRun(pattern: string, signal?: AbortSignal) {
    parseCronExpression(pattern);
    const client = await getSourceServiceClient({ signal });
    return client.getCronNextTime({
      cronExpression: pattern,
    });
  }

  const { value = $bindable([]) }: Props = $props();

  type ScheduleEntry = {
    pattern: string;
    nextRun?: Date;
    error?: string;
    zone?: string;
  };

  const schedules: Promise<ScheduleEntry[]> = $derived(
    Promise.allSettled(value.map((pattern) => getNextRun(pattern))).then(
      (results) =>
        results.map((res, idx) =>
          res.status === "fulfilled"
            ? {
                pattern: value[idx],
                nextRun: toDate(res.value.nextTime!),
                zone: res.value.zone || "UTC",
              }
            : {
                pattern: value[idx],
                error:
                  res.reason instanceof Error
                    ? res.reason.message
                    : `${res.reason}`,
              },
        ),
    ),
  );

  type PatternOption = {
    value: string;
    label: string;
  };

  function genDropdownDays(): PatternOption[] {
    const hours = Array.from({ length: 24 }).keys();
    const dayOfWeeks = [
      "Sunday",
      "Monday",
      "Tuesday",
      "Wednesday",
      "Thursday",
      "Friday",
      "Saturday",
    ];
    const patterns: PatternOption[] = [];
    for (const hour of hours) {
      for (const day of dayOfWeeks) {
        const d = day.slice(0, 3).toUpperCase();
        patterns.push({
          value: `0 ${hour} * * ${d}`,
          label: `At ${hour}:00 on ${day}`,
        });
      }
      patterns.push({
        value: `0 ${hour} * * *`,
        label: `At ${hour}:00 every day`,
      });
    }
    return patterns;
  }

  const patternDropdownList: PatternOption[] = [
    {
      value: "* * * * *",
      label: "Runs every minute",
    },
    {
      value: "@minutely",
      label: "Runs every minute",
    },
    {
      value: "*/5 * * * *",
      label: "Runs every 5 minutes",
    },
    {
      value: "*/15 * * * *",
      label: "Runs every 15 minutes",
    },
    {
      value: "*/30 * * * *",
      label: "Runs every 30 minutes",
    },
    {
      value: "0 * * * *",
      label: "Runs at the start of every hour",
    },
  ].concat(genDropdownDays());

  const filteredPatternDropdownList = $derived.by(() => {
    const text = scheduleInputValue.trim().toLowerCase();
    if (!text) {
      return patternDropdownList;
    }
    return patternDropdownList.filter(
      (opt) =>
        opt.value.toLowerCase().includes(text) ||
        opt.label.toLowerCase().includes(text),
    );
  });

  let showDropdown = $state(false);

  let inputField: HTMLInputElement | null = $state(null);

  const isMobile = window.innerWidth < 640;
</script>

<Combobox.Root
  type="single"
  inputValue={scheduleInputValue}
  onValueChange={(value) => (scheduleInputValue = value.trim())}
  open={showDropdown}
>
  <fieldset class="fieldset">
    <legend class="fieldset-legend">
      <span>Schedule</span>
    </legend>
    <Combobox.Input
      placeholder="e.g. 0 0 * * FRI (Every midnight at Friday)"
      oninput={(e) => (scheduleInputValue = e.currentTarget.value)}
      onfocus={() => (showDropdown = true)}
      onblur={() => (showDropdown = false)}
      class="input w-full"
      bind:ref={inputField}
    ></Combobox.Input>
  </fieldset>
  <Combobox.Portal>
    <Combobox.Content data-theme={$theme} class="z-[9999] bg-base-100">
      <Combobox.Viewport
        class="menu bg-base-100 border-base-content shadow-md"
        style={{
          width: inputField
            ? `${inputField.getBoundingClientRect().width}px`
            : undefined,
        }}
      >
        {#each filteredPatternDropdownList as { value, label }, i (value)}
          {#if i > 0}
            <div class="divider my-0"></div>
          {/if}
          <Combobox.Item
            {value}
            label={value}
            class="btn btn-ghost justify-between"
          >
            {#if !isMobile}
              <span>{value}</span>
            {/if}
            <span class="font-normal">{label}</span>
          </Combobox.Item>
        {:else}
          <Combobox.Item
            label={scheduleInputValue}
            value={scheduleInputValue}
            class="btn btn-ghost justify-between"
          >
            {#if !isMobile}
              <span>{scheduleInputValue}</span>
            {/if}
            <span class="font-normal">Custom Expression</span>
          </Combobox.Item>
        {/each}
      </Combobox.Viewport>
    </Combobox.Content>
  </Combobox.Portal>
</Combobox.Root>

{#snippet helpText()}
  <Popover.Root>
    <Popover.Trigger class="btn btn-square btn-ghost btn-xs">
      <IconInfo />
    </Popover.Trigger>
    <Popover.Portal>
      <Popover.Content class="z-[9999] bg-transparent" data-theme="dracula"
      ></Popover.Content>
    </Popover.Portal>
  </Popover.Root>
{/snippet}
