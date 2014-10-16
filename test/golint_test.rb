require 'test_helper'

class TestGolint < MiniTest::Unit::TestCase
  def test_server
    file = File.open('test/fixtures/server.go')
    results = Golint.lint(file.read)

    assert_equal results.size, 11

    assert_equal results[0].line, 30
    assert_equal results[0].comment, "don't use underscores in Go names; var item_id should be itemID"
    assert_equal results[1].line, 32
    assert_equal results[1].comment, "don't use underscores in Go names; var some_item should be someItem"
    assert_equal results[2].line, 69
    assert_equal results[2].comment, "don't use underscores in Go names; var revision_file should be revisionFile"
    assert_equal results[3].line, 83
    assert_equal results[3].comment, "don't use underscores in Go names; var category_id should be categoryID"
    assert_equal results[4].line, 84
    assert_equal results[4].comment, "don't use underscores in Go names; var aspect_filters should be aspectFilters"
    assert_equal results[5].line, 87
    assert_equal results[5].comment, "don't use underscores in Go names; var per_page should be perPage"
    assert_equal results[6].line, 92
    assert_equal results[6].comment, "don't use underscores in Go names; var page_num should be pageNum"
    assert_equal results[7].line, 97
    assert_equal results[7].comment, "don't use underscores in Go names; var search_result should be searchResult"
    assert_equal results[8].line, 118
    assert_equal results[8].comment, "don't use underscores in Go names; var secret_useragent should be secretUseragent"
    assert_equal results[9].line, 122
    assert_equal results[9].comment, "don't use ALL_CAPS in Go names; use CamelCase"
    assert_equal results[10].line, 122
    assert_equal results[10].comment, "exported const ROUTE_PREFIX should have comment (or a comment on this block) or be unexported"
  end

  def test_one_line
    results = Golint.lint('item_id := vars["item_id"]')

    assert_equal results.size, 1
    assert_equal results[0].line, 1
    assert_equal results[0].comment, "expected 'package', found 'IDENT' item_id"
  end

  def test_valid
    results = Golint.lint("package main\nvar fooBar string")

    assert_equal results.size, 0
  end
end
