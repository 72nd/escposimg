#import "lib.typ": tpl

// ================================================
// CONFIG
// ================================================

#let min_size = 4
#let max_size = 15
#let fonts = ("Arial", "Courier New", "Times New Roman", "Verdana", "Palatino", "Comic Sans MS", "Impact")

#let short_text = "Hamburgefonts"
#let spread_test = "HillimillihirtzheftpflasterEntferner"
#let pangram = "Nyx’ Bö drückt Vamps Quiz-Floß jäh weg."


// ================================================
// LAYOUT & STYLE
// ================================================

#show: tpl.with()

// ================================================
// METHODS
// ================================================

#let measure_text(font_family, font_size, content) = context [
  #calc.round(measure(text(font: font_family, size: font_size, content)).width.mm(), digits: 1)#sym.space.narrow\mm
]

#let entry(font_family, size, content) = {
  let font_size = size * 1pt
  (
    table.cell()[#sym.ballot],
    [
      #font_family \@ #size#sym.space.narrow\pt #sym.arrow #measure_text(font_family, font_size, content)\
      #set text(font: font_family, size: font_size); #content
    ],
  )
}

#let entry_table(entries) = table(
  columns: (auto, 1fr),
  inset: (
    top: 0.0mm,
    bottom: 3.0mm,
    left: 0.4mm,
    right: 0.4mm,
  ),
  stroke: none,
  ..entries
)

#let rndr_size_range(font_family, content) = for size in range(min_size, max_size) {
  entry(font_family, size, content)
}

#let rndr_all_fonts(font_size, content) = for font in fonts {
  entry(font, font_size, content)
}

#let rndr_all_fonts_with_variation(font_size, content) = for font in fonts {
  entry(font, font_size, [#content\ #emph(content)\ #strong(content)])
}

// ================================================
// CONTENT
// ================================================


#align(center)[*Fonts*]

#entry_table(
  rndr_all_fonts_with_variation(8, pangram)
)
