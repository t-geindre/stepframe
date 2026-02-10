# Stepframe — A live-oriented MIDI sequencer — Feature Roadmap

## V1 — Playable MVP (Live-first, stable timing, minimal editing)

### Global Transport & Tempo
- **Global Play / Stop**  
  Starts and stops the whole engine. When stopped, no new events are scheduled or sent.
- **Tempo (BPM)**  
  Sets the global tempo used to compute timing. Must be adjustable while stopped (and ideally while playing).
- **Basic Time Signature (default 4/4)**  
  Defines bar boundaries for quantized launching. V1 can hardcode 4/4 but should represent the concept internally.

### MIDI Clock Sync (Master / Slave)
- **Clock Mode: Internal (Master)**  
  Sequencer generates timing internally and can optionally send MIDI Clock + Start/Stop to an output.
- **Clock Mode: External (Slave)**  
  Sequencer follows incoming MIDI Clock and responds to MIDI Start/Stop. If clock stops, playback must stop or clearly indicate “out of sync”.
- **Send MIDI Clock (optional toggle)**  
  When enabled in Master mode, sends MIDI Clock ticks + Start/Stop messages to selected output port.

### Tracks (Core Live Workflow)
- **Multiple MIDI Tracks**  
  A project contains N tracks (configurable or fixed count in V1). Each track holds a MIDI clip/pattern.
- **Per-Track Play / Stop (Track Enable)**  
  Each track can be enabled/disabled while the global engine runs. Disabling a track stops it from emitting new events (and should send “All Notes Off” for that track/channel if needed).
- **Launch Quantization (Bar or Beat)**  
  When enabling a track, its start aligns to the next bar or beat (user-selectable per track or global in V1). No mid-beat starts unless explicitly allowed.
- **MIDI Output Assignment (Port + Channel)**  
  Each track routes to a selected MIDI output port and a MIDI channel (1–16). All events produced by that track are sent there.

### Track Timing / Length
- **Track Length**  
  Each track has a loop length (in bars, beats, or steps). Playback wraps cleanly at the loop end.
- **Per-Track Rate / Division**  
  Each track can run at a musical division relative to the global clock (e.g. 1/1, 1/2, 1/4…). This enables simple polyrhythmic layering.

### Recording (From MIDI Input)
- **Select MIDI Input Source**  
  Choose which input port(s) are recorded (e.g. a master keyboard).
- **Record Arm per Track**  
  Only armed tracks record incoming notes.
- **Real-time Recording into a Track**  
  Incoming MIDI Note On/Off events are captured into the armed track with timestamps.
- **Overdub vs Replace (simple)**  
  Defines recording behavior. Overdub adds notes without deleting existing ones; replace clears the loop before recording.
- **Basic Record Quantize (On/Off)**  
  Optional quantization during recording (e.g. snap note start times to nearest grid). Keep it simple in V1.

### Editing & Visualization (Minimal)
- **Basic Piano Roll Viewer**  
  A single-track view that displays notes (pitch vs time) for verification.
- **Basic Note Editing**  
  Add/delete notes and move them in time. Duration editing can be minimal in V1.

### Safety / Reliability
- **MIDI Panic / All Notes Off**  
  A dedicated action that sends All Notes Off (CC123) and/or Note Off for all active notes.
- **Stuck Note Prevention**  
  Ensure every Note On has a corresponding Note Off even when stopping, muting a track, or changing quantization.

### Persistence
- **Save / Load Project**  
  Persist tracks, clips, routing (port/channel), tempo, and quantization settings to disk.

---

## V2 — Serious Live Tool (Better editing, better sync, better control)

### Advanced Sync / Timing
- **Launch Quantization: Step / Beat / Bar**  
  More granular launch options, including step-level quantization.
- **MIDI Song Position Pointer (SPP) Support**  
  When following external clock, respond to SPP to align playback position.
- **Latency / Clock Offset Compensation**  
  User-configurable timing offsets to compensate for device latency.

### Track Controls
- **Mute vs Stop Behavior**  
  Mute silences output but keeps the playhead running; Stop disables and resets or holds position.
- **Solo**  
  Solo one track while muting others.
- **Per-Track All Notes Off on State Change**  
  Optionally send All Notes Off when muting or stopping a track.

### Editing Upgrades (Piano Roll)
- **Velocity Editing**  
  Edit per-note velocity values.
- **Multi-Select & Batch Operations**  
  Move, transpose, or resize multiple notes at once.
- **Zoom & Scroll**  
  Zoom and scroll in time and pitch, suitable for small touch displays.

### Recording Upgrades
- **Non-Destructive Quantization**  
  Apply quantization without destroying original timing data.
- **Pre-roll / Count-in**  
  Optional count-in before recording begins.
- **Undo / Redo**  
  Undo and redo edits and recording actions.

### MIDI Data Beyond Notes
- **Control Change (CC) Recording & Playback**  
  Record and play back MIDI CC messages.
- **Program Change**  
  Send program change messages per track.

### Project / File Handling
- **Export Standard MIDI File (SMF)**  
  Export tracks or clips as MIDI files.
- **Import MIDI File (basic)**  
  Import MIDI files into tracks.

### UX / Live Ergonomics
- **MIDI Learn**  
  Map external controller inputs to sequencer functions.
- **Performance Screen**  
  A dedicated live-performance-oriented view.

---

## V3 — Pro Features (Polish, depth, advanced workflows)

### Advanced Track & Rhythm Features
- **Independent Track Loop Length (Non-multiples)**  
  Allow odd loop lengths across tracks.
- **Per-Track Swing / Groove**  
  Apply swing or groove templates per track.
- **Track Phase / Delay**  
  Shift track timing forward or backward for groove alignment.

### Automation / Expressive MIDI
- **Pitch Bend Support**  
  Record, edit, and play pitch bend data.
- **Aftertouch (Channel / Poly)**  
  Support channel and polyphonic aftertouch.
- **CC Lanes Editing**  
  Graphical editing of CC automation lanes.

### Live Performance Workflows
- **Scenes**  
  Recall sets of track states with quantized launching.
- **Clip / Pattern Variations**  
  Multiple clips per track for live arrangement.
- **Follow Actions (optional)**  
  Automatically trigger clips after a defined duration.

### Robustness & Observability
- **Sync Status & Diagnostics View**  
  Visual indicators for clock lock and timing health.
- **Event Logging / MIDI Monitor**  
  Inspect incoming and outgoing MIDI events.

### Power User Features
- **Templates & Default Routing**  
  Reusable project templates and default MIDI routing.
- **Project Versioning / Snapshots**  
  Save and recall snapshots of a project state.
