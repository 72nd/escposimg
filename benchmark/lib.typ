#let tpl(
  height: auto,
  font_size: 8pt,
  right_margin: 7.5mm,
  body,
) = {
  set page(width: 79.5mm, height: height, margin: (
    top: 0mm,
    left: 0mm,
    right: right_margin,
    bottom: 0mm,
  ))

  set par(
    spacing: 0.9em,
    leading: 0.4em,
  )

  show heading.where(level: 2): set text(size: 0.9em)
  show heading.where(level: 2): set block(above: .9em)

  set text(
    size: font_size,
    font: "Arial",
    hyphenate: true,
  )

  show raw: set text(font: "Courier New", size: 10pt)

  show heading.where(level: 1): set text(size: 0.6em)

  body
}

// Measurement tape scale with cm and mm markings
#let measure_tape(
  width: 79.5mm,
  cm_count: 8,
) = {
  // Main ruler line
  line(length: width, stroke: 0.5pt)

  // Generate cm and mm marks
  for cm in range(cm_count) {
    let x_pos = cm * 10mm

    // CM mark (longer line)
    place(top + left, dx: x_pos, dy: 0mm, line(angle: 90deg, length: 3mm, stroke: 0.5pt))

    // CM number
    place(top + left, dx: x_pos + 1mm, dy: 3.5mm, text(size: 6pt, str(cm)))

    // MM marks (shorter lines)
    for mm in range(1, 10) {
      let mm_x = x_pos + mm * 1mm
      if mm_x < width {
        let mark_height = if mm == 5 { 2mm } else { 1mm }
        place(top + left, dx: mm_x, dy: 0mm, line(angle: 90deg, length: mark_height, stroke: 0.3pt))
      }
    }
  }
}
