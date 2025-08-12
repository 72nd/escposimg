#import "lib.typ": measure_tape, tpl

// ================================================
// LAYOUT & STYLE
// ================================================

#let page_height = 25mm

#show: tpl.with(
  height: page_height,
)


// ================================================
// METHODS
// ================================================

#let vertical_measure_tape() = {
  let cm_count = 3
  let cutoff = 7mm

  for cm in range(cm_count) {
    let y_pos = cm * 10mm
    place(top + left, dy: y_pos, line(length: 3mm, stroke: 0.5pt))
    if cm != 0 {
      place(top + left, dx: 4mm, dy: y_pos, text(size: 6pt, str(cm)))
    }

    for mm in range(1, 10) {
      let y_mm = y_pos + mm * 1mm
      if y_mm < cutoff {
        continue
      }
      let mark_height = if mm == 5 { 2mm } else { 1mm }
      place(top + left, dy: y_mm, line(length: mark_height, stroke: 0.3pt))
    }
  }
}


// ================================================
// CONTENT
// ================================================

#measure_tape()

#vertical_measure_tape()

#align(center + horizon)[
  This layout is 79.5#sym.space.narrow\mm by #page_height.mm()#sym.space.narrow\mm.
]

