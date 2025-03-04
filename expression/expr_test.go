/*
 * TencentBlueKing is pleased to support the open source community by making 蓝鲸智云PaaS平台社区版 (BlueKing PaaS
 * Community Edition) available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package expression_test

import (
	"strconv"
	"testing"

	"github.com/TencentBlueKing/iam-go-sdk/expression"
	"github.com/TencentBlueKing/iam-go-sdk/expression/operator"

	"github.com/stretchr/testify/assert"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Expr", func() {

	Describe("Eval", func() {
		var e *expression.ExprCell
		var o expression.ObjectSetInterface
		BeforeEach(func() {
			o = expression.NewObjectSet()
		})

		It("op.AND", func() {
			e = &expression.ExprCell{
				OP: operator.AND,
				Content: []expression.ExprCell{
					{
						OP:    operator.Eq,
						Field: "obj.id",
						Value: 1,
					},
					{
						OP:    operator.Eq,
						Field: "obj.name",
						Value: "object",
					},
				},
			}

			// String
			assert.Equal(GinkgoT(), "((obj.id eq 1) AND (obj.name eq object))", e.String())

			//hit
			o.Set("obj", map[string]interface{}{
				"id":   1,
				"name": "object",
			})
			assert.True(GinkgoT(), e.Eval(o))

			// Render
			assert.Equal(GinkgoT(), "((1 eq 1) AND (object eq object))", e.Render(o))

			//miss
			o.Set("obj", map[string]interface{}{
				"id":   2,
				"name": "object",
			})
			assert.False(GinkgoT(), e.Eval(o))

		})

		It("op.OR", func() {
			e = &expression.ExprCell{
				OP: operator.OR,
				Content: []expression.ExprCell{
					{
						OP:    operator.Eq,
						Field: "obj.id",
						Value: 1,
					},
					{
						OP:    operator.Eq,
						Field: "obj.name",
						Value: "object",
					},
				},
			}

			// String
			assert.Equal(GinkgoT(), "((obj.id eq 1) OR (obj.name eq object))", e.String())

			//hit
			o.Set("obj", map[string]interface{}{
				"id":   1,
				"name": "object1",
			})
			assert.True(GinkgoT(), e.Eval(o))

			// Render
			assert.Equal(GinkgoT(), "((1 eq 1) OR (object1 eq object))", e.Render(o))

			//miss
			o.Set("obj", map[string]interface{}{
				"id":   2,
				"name": "object2",
			})
			assert.False(GinkgoT(), e.Eval(o))

		})

		Context("op.BinaryOperator", func() {

			It("op.Any", func() {
				e = &expression.ExprCell{
					OP:    operator.Any,
					Field: "obj.id",
					Value: nil,
				}
				assert.True(GinkgoT(), e.Eval(o))
			})

			Context("evalPositive", func() {
				It("op.Eq", func() {
					e = &expression.ExprCell{
						OP:    operator.Eq,
						Field: "obj.name",
						Value: "hello",
					}
					// true
					o.Set("obj", map[string]interface{}{
						"name": "hello",
					})
					assert.True(GinkgoT(), e.Eval(o))

					// false
					o.Set("obj", map[string]interface{}{
						"name": "abc",
					})
					assert.False(GinkgoT(), e.Eval(o))

					// type not match
					o.Set("obj", map[string]interface{}{
						"name": 1,
					})
					assert.False(GinkgoT(), e.Eval(o))
				})

				It("op.Eq value is an array", func() {
					e = &expression.ExprCell{
						OP:    operator.Eq,
						Field: "obj.name",
						Value: "hello",
					}

					// hit
					o.Set("obj", map[string]interface{}{
						"name": []string{"hello", "world"},
					})
					assert.True(GinkgoT(), e.Eval(o))

					// miss
					o.Set("obj", map[string]interface{}{
						"name": []string{"abc", "def"},
					})
					assert.False(GinkgoT(), e.Eval(o))
				})

				// lt/lte/gt/gte/starts_with/ends_with/in
				Context("lt/lte/gt/gte/starts_with/ends_with/in", func() {
					It("lt", func() {
						e = &expression.ExprCell{
							OP:    operator.Lt,
							Field: "obj.age",
							Value: 18,
						}

						// hit
						o.Set("obj", map[string]interface{}{
							"age": 17,
						})
						assert.True(GinkgoT(), e.Eval(o))

						// miss
						o.Set("obj", map[string]interface{}{
							"age": 18,
						})
						assert.False(GinkgoT(), e.Eval(o))
					})
					It("lte", func() {
						e = &expression.ExprCell{
							OP:    operator.Lte,
							Field: "obj.age",
							Value: 18,
						}

						// hit
						o.Set("obj", map[string]interface{}{
							"age": 18,
						})
						assert.True(GinkgoT(), e.Eval(o))

						// miss
						o.Set("obj", map[string]interface{}{
							"age": 19,
						})
						assert.False(GinkgoT(), e.Eval(o))
					})

					It("gt", func() {
						e = &expression.ExprCell{
							OP:    operator.Gt,
							Field: "obj.age",
							Value: 18,
						}

						// hit
						o.Set("obj", map[string]interface{}{
							"age": 19,
						})
						assert.True(GinkgoT(), e.Eval(o))

						// miss
						o.Set("obj", map[string]interface{}{
							"age": 18,
						})
						assert.False(GinkgoT(), e.Eval(o))

					})
					It("gte", func() {
						e = &expression.ExprCell{
							OP:    operator.Gte,
							Field: "obj.age",
							Value: 18,
						}

						// hit
						o.Set("obj", map[string]interface{}{
							"age": 18,
						})
						assert.True(GinkgoT(), e.Eval(o))

						// miss
						o.Set("obj", map[string]interface{}{
							"age": 17,
						})
						assert.False(GinkgoT(), e.Eval(o))

					})

					It("starts_with", func() {
						e = &expression.ExprCell{
							OP:    operator.StartsWith,
							Field: "obj.name",
							Value: "hello",
						}

						// hit
						o.Set("obj", map[string]interface{}{
							"name": "hello world",
						})
						assert.True(GinkgoT(), e.Eval(o))

						// miss
						o.Set("obj", map[string]interface{}{
							"name": "foo bar",
						})
						assert.False(GinkgoT(), e.Eval(o))

						// NOTE: _bk_iam_path_ with starts_with
					})
					It("starts_with with _bk_iam_path_", func() {
						e = &expression.ExprCell{
							OP:    operator.StartsWith,
							Field: "obj._bk_iam_path_",
							Value: "/a,1/b,*/",
						}

						o.Set("obj", map[string]interface{}{
							"_bk_iam_path_": "/a,1/b,2/c,3/",
						})
						assert.True(GinkgoT(), e.Eval(o))
					})

					It("ends_with", func() {
						e = &expression.ExprCell{
							OP:    operator.EndsWith,
							Field: "obj.name",
							Value: "hello",
						}

						// hit
						o.Set("obj", map[string]interface{}{
							"name": "world hello",
						})
						assert.True(GinkgoT(), e.Eval(o))

						// miss
						o.Set("obj", map[string]interface{}{
							"name": "foo bar",
						})
						assert.False(GinkgoT(), e.Eval(o))
					})
					It("in", func() {
						e = &expression.ExprCell{
							OP:    operator.In,
							Field: "obj.name",
							Value: []string{"hello", "world"},
						}

						// hit
						o.Set("obj", map[string]interface{}{
							"name": "hello",
						})
						assert.True(GinkgoT(), e.Eval(o))

						// miss
						o.Set("obj", map[string]interface{}{
							"name": "foo",
						})
						assert.False(GinkgoT(), e.Eval(o))
					})
				})
			})

			Context("evalNegative", func() {
				It("op.NotEq", func() {
					e = &expression.ExprCell{
						OP:    operator.NotEq,
						Field: "obj.name",
						Value: "hello",
					}
					o.Set("obj", map[string]interface{}{
						"name": "world",
					})
					assert.True(GinkgoT(), e.Eval(o))
				})

				It("op.NotEq value is an array", func() {
					e = &expression.ExprCell{
						OP:    operator.NotEq,
						Field: "obj.name",
						Value: "hello",
					}

					// false, all not eq
					o.Set("obj", map[string]interface{}{
						"name": []string{"abc", "def"},
					})
					assert.True(GinkgoT(), e.Eval(o))

					// true, one equal
					o.Set("obj", map[string]interface{}{
						"name": []string{"hello", "world"},
					})
					assert.False(GinkgoT(), e.Eval(o))
				})

				// not_starts_with/not_ends_with/not_in
				Context("not_starts_with/not_ends_with/not_in", func() {
					It("not_starts_with", func() {
						e = &expression.ExprCell{
							OP:    operator.NotStartsWith,
							Field: "obj.name",
							Value: "hello",
						}

						// hit
						o.Set("obj", map[string]interface{}{
							"name": "foo bar",
						})
						assert.True(GinkgoT(), e.Eval(o))

						// miss
						o.Set("obj", map[string]interface{}{
							"name": "hello world",
						})
						assert.False(GinkgoT(), e.Eval(o))
					})
					It("not_ends_with", func() {
						e = &expression.ExprCell{
							OP:    operator.NotEndsWith,
							Field: "obj.name",
							Value: "hello",
						}

						// hit
						o.Set("obj", map[string]interface{}{
							"name": "foo bar",
						})
						assert.True(GinkgoT(), e.Eval(o))

						// miss
						o.Set("obj", map[string]interface{}{
							"name": "world hello",
						})
						assert.False(GinkgoT(), e.Eval(o))

					})
					It("not_in", func() {
						e = &expression.ExprCell{
							OP:    operator.NotIn,
							Field: "obj.name",
							Value: []string{"hello", "world"},
						}

						// hit
						o.Set("obj", map[string]interface{}{
							"name": "foo",
						})
						assert.True(GinkgoT(), e.Eval(o))

						// miss
						o.Set("obj", map[string]interface{}{
							"name": "hello",
						})
						assert.False(GinkgoT(), e.Eval(o))

					})
				})
			})

			It("op.Contains", func() {
				e = &expression.ExprCell{
					OP:    operator.Contains,
					Field: "obj.name",
					Value: "hello",
				}
				o.Set("obj", map[string]interface{}{
					"name": []string{"hello", "world"},
				})
				assert.True(GinkgoT(), e.Eval(o))
			})

			It("op.NotContains", func() {
				e = &expression.ExprCell{
					OP:    operator.NotContains,
					Field: "obj.name",
					Value: "abc",
				}
				o.Set("obj", map[string]interface{}{
					"name": []string{"hello", "world"},
				})
				assert.True(GinkgoT(), e.Eval(o))
			})

		})

	})

})

func BenchmarkExprCellEqual(b *testing.B) {
	e := &expression.ExprCell{
		OP:    operator.Eq,
		Field: "obj.name",
		Value: "hello",
	}

	o := expression.NewObjectSet()
	o.Set("obj", map[string]interface{}{
		"name": "world",
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Eval(o)
	}
}
func BenchmarkExprCellNotEqual(b *testing.B) {
	e := &expression.ExprCell{
		OP:    operator.NotEq,
		Field: "obj.name",
		Value: "hello",
	}

	o := expression.NewObjectSet()
	o.Set("obj", map[string]interface{}{
		"name": "hello",
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Eval(o)
	}
}
func BenchmarkExprCellLess(b *testing.B) {
	e := &expression.ExprCell{
		OP:    operator.Lt,
		Field: "obj.age",
		Value: 18,
	}

	o := expression.NewObjectSet()
	o.Set("obj", map[string]interface{}{
		"name": "helloworld",
		"age":  2,
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Eval(o)
	}
}
func BenchmarkExprCellStartsWith(b *testing.B) {
	e := &expression.ExprCell{
		OP:    operator.StartsWith,
		Field: "obj.name",
		Value: "hello",
	}

	o := expression.NewObjectSet()
	o.Set("obj", map[string]interface{}{
		"name": "helloworld",
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Eval(o)
	}
}
func BenchmarkExprCellIn(b *testing.B) {

	ids := make([]string, 10000)
	for i := 0; i < 9999; i++ {
		ids = append(ids, strconv.Itoa(i))
	}
	ids = append(ids, "world")

	e := &expression.ExprCell{
		OP:    operator.In,
		Field: "obj.name",
		// Value: []string{"hello", "world"},
		Value: ids,
	}

	o := expression.NewObjectSet()
	o.Set("obj", map[string]interface{}{
		"name": "world",
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Eval(o)
	}
}
