/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *modelClassification) isExists(ctx core.ContextParams, classificationID string) (origin *metadata.Classification, exists bool, err error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.ClassFieldClassificationID, Val: ctx.SupplierAccount}, &mongo.Eq{Key: metadata.ClassFieldClassificationID, Val: classificationID})
	err = m.dbProxy.Table(common.BKTableNameObjClassifiction).Find(cond.ToMapStr()).One(ctx, origin)
	return origin, m.dbProxy.IsNotFoundError(err), err
}

func (m *modelClassification) hasModel(ctx core.ContextParams, cond universalsql.Condition) (cnt uint64, exists bool, err error) {

	cnt, err = m.dbProxy.Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).Count(ctx)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database count operation on the table(%s) by the condition(%v), error info is %s", ctx.ReqID, common.BKTableNameObjDes, cond.ToMapStr(), err.Error())
		return 0, false, err
	}
	exists = 0 != cnt
	return cnt, exists, err
}

func (m *modelClassification) cascadeDeleteModel(ctx core.ContextParams, classificationIDS []string) (uint64, error) {

	deleteCond := mongo.NewCondition()
	deleteCond.Element(&mongo.In{Key: metadata.ModelFieldObjCls, Val: classificationIDS})
	deleteCond.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})
	return m.model.cascadeDelete(ctx, deleteCond)
}